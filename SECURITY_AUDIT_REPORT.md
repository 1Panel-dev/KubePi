# KubePi 安全审计报告

审计日期：2026-06-22

## 审计范围

- 后端 Go 代码：`cmd/`、`internal/`、`pkg/`、`migrate/`、`thirdparty/gotty/`
- 前端代码与依赖：`web/dashboard/`、`web/kubepi/`、`web/terminal/`、`thirdparty/gotty/js/`
- 重点关注：认证与授权、Kubernetes 高权限操作、敏感信息泄露、命令注入、SSRF/TLS、会话安全、依赖漏洞。

## 总体结论

本次审计确认存在多处高风险问题。最核心的风险集中在两类：

1. KubePi 自身的 RBAC 和集群内 Kubernetes RBAC 没有在若干高危接口上形成一致的权限边界，普通已登录用户可能借助 KubePi 的集群导入凭据执行超出自身权限的 Pod exec、日志读取、Helm 安装/卸载、Pod 文件读写等操作。
2. 多个列表/详情接口直接返回包含凭据的持久化模型，导致只读角色也可能读取集群 kubeconfig、bearer token、私钥、LDAP 密码、镜像仓库密码、监控 token、用户 MFA secret 等敏感信息。

建议优先处理 `KBP-SEC-001` 至 `KBP-SEC-006`。这些问题相互叠加后，影响可能从普通账号扩大到集群接管、宿主机接管或平台账号接管。

## 风险概览

| ID | 严重级别 | 问题 |
| --- | --- | --- |
| KBP-SEC-001 | 严重 | 集群列表/详情接口泄露集群凭据和私钥 |
| KBP-SEC-002 | 严重 | `charts` / `apps` / `pod` 被加入白名单，绕过通用 RBAC |
| KBP-SEC-003 | 严重 | 终端、节点终端使用集群导入凭据，普通用户可越权执行 Pod/节点 Shell |
| KBP-SEC-004 | 严重 | Pod 文件管理接口越权读写容器文件，并存在命令注入风险 |
| KBP-SEC-005 | 高 | SSO 自动创建的本地账号使用固定默认密码 |
| KBP-SEC-006 | 高 | 多类敏感配置和 MFA secret 通过只读接口泄露 |
| KBP-SEC-007 | 高 | 日志流接口使用集群导入凭据，可越权读取任意 Pod 日志 |
| KBP-SEC-008 | 高 | SAML 依赖存在可达签名绕过漏洞 |
| KBP-SEC-009 | 中 | SSO 流程使用全局可变状态、固定 state，并信任 Referer 生成回调地址 |
| KBP-SEC-010 | 中 | 多处关闭 TLS 证书校验 |
| KBP-SEC-011 | 中 | Cookie 会话缺少显式安全属性和 CSRF 防护 |
| KBP-SEC-012 | 中 | 前端生产依赖存在多个 npm audit 告警 |

## 详细发现

### KBP-SEC-001 集群列表/详情接口泄露集群凭据和私钥

严重级别：严重

影响：

- 拥有 `clusters get/list` 权限的用户可以通过集群详情、列表或搜索接口读取完整 `Cluster` 模型。
- `Cluster` 模型包含 `privateKey`、`bearerToken`、`configFileContent`、证书私钥、代理密码等字段。
- 默认 `Common User` 具备 `clusters get/list`，默认 `ReadOnly` 具备 `* get/list`。因此该问题不是管理员专属路径。
- 泄露的 kubeconfig、token 或证书私钥可被用于直接访问 Kubernetes API，绕过 KubePi 的应用层权限控制。

证据：

- 敏感字段定义：`internal/model/v1/cluster/cluster.go:7-48`
- `GetCluster` 原样返回 `c`：`internal/api/v1/cluster/cluster.go:411-420`
- `ListClusters` 原样嵌入 `Cluster: clusters[i]` 并返回：`internal/api/v1/cluster/cluster.go:500-530`
- `SearchClusters` 原样嵌入 `Cluster: clusters[i]` 并返回：`internal/api/v1/cluster/cluster.go:288-315`
- 默认角色：`ReadOnly` 为 `* get/list`，`Common User` 为 `clusters get/list`：`migrate/v1/migrations.go:70-109`

建议：

- 为所有集群响应建立专用 DTO，默认不返回 `PrivateKey`、`BearerToken`、`Certificate.KeyData`、`ConfigFileContent`、代理密码等字段。
- 只有管理员在必要的专用接口中才能查看或轮换敏感凭据，且返回前应做掩码。
- 非管理员查询集群时，应先校验 `clusterBinding`，并只返回已绑定集群的非敏感元信息。
- 增加回归测试：普通 `Common User` / `ReadOnly` 调用集群列表、详情、搜索接口时响应中不得出现任何 secret 字段。

### KBP-SEC-002 `charts` / `apps` / `pod` 被加入白名单，绕过通用 RBAC

严重级别：严重

影响：

- `resourceWhiteList` 中包含 `charts`、`apps`、`pod`，`roleAccessHandler` 对这些资源跳过通用 RBAC。
- 这些资源不是纯静态或低风险资源，包含 Helm 仓库增删改、Chart 安装/升级/卸载、Pod 文件读写等真实集群变更能力。
- 任意已登录用户只要通过 `authHandler`，即可进入这些业务接口，具体权限依赖后续 handler 是否自行校验；当前多个 handler 没有做用户绑定或管理员校验。

证据：

- 白名单定义：`internal/api/v1/v1.go:46`
- 白名单资源跳过 RBAC：`internal/api/v1/v1.go:339-345`
- `charts`、`apps` 路由包含写操作：`internal/api/v1/chart/chart.go:272-291`
- `pod` 路由包含文件读写、上传、下载、删除：`internal/api/v1/file/file.go:275-287`

建议：

- 从全局白名单中移除 `charts`、`apps`、`pod`。
- 对确实需要公开的轻量资源拆分独立路径，而不是按一级资源整体放行。
- 对每个 handler 做显式权限要求：KubePi 资源权限、集群绑定关系、Kubernetes SAR/SSAR 权限应同时成立。

### KBP-SEC-003 终端、节点终端使用集群导入凭据，普通用户可越权执行 Pod/节点 Shell

严重级别：严重

影响：

- `GET /clusters/:name/terminal/session` 和 `GET /clusters/:name/node_terminal/session` 属于 `clusters:get` 路由，默认 `Common User` 和 `ReadOnly` 都可能满足应用层权限。
- handler 直接通过 `clusterService.Get` 读取集群配置，再用 `kubernetes.NewKubernetes(c)` 创建集群导入凭据的 client/config，没有使用当前用户的 `clusterBinding` 证书。
- Pod 终端可在任意 namespace/pod/container 中执行 Shell。
- 节点终端会在目标节点创建特权 Pod，开启 `Privileged`、`HostNetwork`、`HostPID`，并将宿主机 `/` 挂载到 `/host`，风险等同于宿主机级别接管。

证据：

- 路由：`internal/api/v1/cluster/cluster.go:625-628`
- Pod 终端使用集群导入凭据：`internal/api/v1/cluster/terminal.go:29-57`
- 节点终端使用集群导入凭据：`internal/api/v1/cluster/terminal.go:74-90`
- 节点特权 Pod：`pkg/terminal/shell.go:335-408`
- 节点 Shell exec：`pkg/terminal/shell.go:419-436`
- 默认角色权限：`migrate/v1/migrations.go:70-109`

建议：

- Pod 终端必须使用当前用户的 `clusterBinding` 证书和私钥，并由 Kubernetes API server 进行 `pods/exec` 授权。
- 节点终端建议限制为管理员专用，并增加二次确认、审计记录和可配置禁用开关。
- 创建特权 Pod 前应强制校验管理员身份和目标集群授权，不应仅依赖 `clusters:get`。

### KBP-SEC-004 Pod 文件管理接口越权读写容器文件，并存在命令注入风险

严重级别：严重

影响：

- `/pod/*` 被放入白名单，普通已登录用户可以进入 Pod 文件管理接口。
- 服务层通过集群导入凭据创建 Kubernetes client，未使用当前用户身份，也没有校验用户是否有目标 namespace/pod/container 的 exec 权限。
- 文件创建和编辑把用户输入的 `content`、`path` 拼接到 `sh -c` 字符串中，存在容器内命令注入风险。
- 这些接口可以读文件、写文件、上传、下载、删除和重命名容器内文件，实际能力接近远程命令执行。

证据：

- `/pod` 路由：`internal/api/v1/file/file.go:275-287`
- 创建文件拼接 `sh -c`：`internal/api/v1/file/file.go:63-73`
- 编辑文件拼接 `sh -c`：`pkg/util/podtool/edit.go:3-17`
- 服务层使用集群导入凭据：`internal/service/v1/file/file.go:48-73`
- Pod exec 入口：`pkg/util/podtool/podtool.go:48-63`
- 下载/上传接口直接接受目标 `cluster/namespace/podName/containerName/path`：`internal/api/v1/file/file.go:156-222`

建议：

- 移除 `/pod` 白名单，所有 Pod 文件接口都必须绑定当前用户身份。
- 使用当前用户 Kubernetes 凭据执行，并让 Kubernetes RBAC 校验 `pods/exec`。
- 避免 `sh -c` 拼接；能用参数数组完成的操作不要经过 shell。
- 对确需 shell 的场景使用严格的路径白名单、shell escaping 和内容流式写入，不把用户输入拼进命令字符串。
- 增加审计日志，记录用户、集群、namespace、pod、container、路径和操作类型。

### KBP-SEC-005 SSO 自动创建的本地账号使用固定默认密码

严重级别：高

影响：

- OpenID/SAML 自动创建用户时，设置了固定密码，并将用户类型设为 `LOCAL`。
- 普通密码登录接口对非 LDAP 用户直接使用本地密码校验。
- 因此只要知道 SSO 自动创建账号的用户名或邮箱，就可能绕过 SSO 入口，用固定密码从普通登录接口登录。
- 自动创建的 SSO 用户默认绑定 `ReadOnly`，结合 `KBP-SEC-001` 后可能进一步读取集群凭据。

证据：

- 普通登录对非 LDAP 用户使用本地密码校验：`internal/api/v1/session/session.go:103-137`
- OpenID 自动创建用户固定密码、类型为 `LOCAL`：`internal/service/v1/sso/sso.go:184-198`
- OpenID 自动绑定 `ReadOnly`：`internal/service/v1/sso/sso.go:212-227`
- SAML 自动创建用户固定密码、类型为 `LOCAL`：`internal/service/v1/sso/sso.go:331-343`
- SAML 自动绑定 `ReadOnly`：`internal/service/v1/sso/sso.go:357-372`

建议：

- SSO 自动创建用户应使用独立类型，例如 `SSO`，并禁止普通密码登录。
- 如果必须落地密码，应生成不可登录的随机高强度值，并要求首次绑定或管理员显式启用本地密码。
- 清理历史自动创建用户的固定密码风险，批量重置或迁移用户类型。

### KBP-SEC-006 多类敏感配置和 MFA secret 通过只读接口泄露

严重级别：高

影响：

- 多个列表/详情接口返回持久化模型，未脱敏密码、token、MFA secret 等字段。
- 默认 `ReadOnly` 角色具备 `* get/list`，可能读取这些敏感字段。
- 泄露的 MFA secret 可削弱 MFA 防护；泄露的 LDAP、镜像仓库、Grafana、Prometheus 凭据可横向访问外部系统。

证据：

- LDAP 列表原样返回，模型包含 `Password`：`internal/api/v1/ldap/ldap.go:20-28`、`internal/model/v1/ldap/ldap.go:8-22`
- 用户接口只清空 `Authenticate`，未清空 `Mfa.Secret`：`internal/api/v1/user/user.go:71-85`、`internal/api/v1/user/user.go:279-290`、`internal/api/v1/user/user.go:311-314`
- 用户模型包含 MFA secret：`internal/model/v1/user/user.go:5-25`
- Grafana/Metric 返回包含 token/password 的模型：`internal/api/v1/monitor/monitor.go:43-48`、`internal/api/v1/monitor/monitor.go:115-123`、`internal/api/v1/monitor/monitor.go:165-173`、`internal/model/v1/monitor/monitor.go:5-31`
- 镜像仓库返回包含 password 的模型：`internal/api/v1/imagerepo/image_repo.go:40-47`、`internal/api/v1/imagerepo/image_repo.go:184-191`、`internal/api/v1/imagerepo/image_repo.go:207-215`、`internal/model/v1/imagerepo/image_repo.go:5-21`

建议：

- 为所有 API 响应建立脱敏 DTO，不直接返回数据库模型。
- 默认隐藏 `password`、`token`、`secret`、`key`、`privateKey`、`configFileContent` 等字段。
- 管理员也只返回掩码，必要时通过专用“重置/轮换/测试连接”流程使用 secret。
- 对用户列表/详情统一清空 `Mfa.Secret`，只返回 `Mfa.Enable` 状态。

### KBP-SEC-007 日志流接口使用集群导入凭据，可越权读取任意 Pod 日志

严重级别：高

影响：

- `GET /clusters/:name/logging/session` 与终端接口一样属于 `clusters:get` 路径。
- handler 使用集群导入凭据读取任意 namespace/pod/container 日志。
- Pod 日志经常包含 token、业务密钥、错误堆栈或用户数据，越权读取会造成敏感信息泄露。

证据：

- 路由：`internal/api/v1/cluster/cluster.go:625-628`
- handler 使用集群导入凭据：`internal/api/v1/cluster/logging.go:61-78`
- 日志流直接调用 Kubernetes `Pods(namespace).GetLogs(pod, ...)`：`pkg/logging/logging.go:121-130`

建议：

- 日志读取必须使用当前用户 Kubernetes 凭据，并让 Kubernetes RBAC 校验 `pods/log`。
- 如果需要跨 namespace 聚合，应先计算用户可访问 namespace，不得直接接受任意 namespace。
- session ID 建议使用 `crypto/rand` 生成，当前 `pkg/logging/logging.go:19-26` 使用 `math/rand`，虽然不是主要越权根因，但应一并加固。

### KBP-SEC-008 SAML 依赖存在可达签名绕过漏洞

严重级别：高

影响：

- `govulncheck` 确认项目代码可达 `github.com/russellhaering/goxmldsig@v1.4.0` 的 `GO-2026-4753`。
- 漏洞标题为 `Loop Variable Capture Signature Bypass in goxmldsig`，修复版本为 `v1.6.0`。
- 该依赖由 `github.com/crewjam/saml@v0.5.1` 引入，且项目 SAML 登录路径可达相关代码。

证据：

- `go.mod` 直接依赖 `github.com/crewjam/saml v0.5.1` 和 `github.com/russellhaering/goxmldsig v1.4.0`。
- `go mod graph` 显示 `github.com/crewjam/saml@v0.5.1 -> github.com/russellhaering/goxmldsig@v1.4.0`。
- `govulncheck` 可达路径包括 `internal/api/v1/sso/sso.go:266` 的 `Saml2Auth` 和 `internal/api/v1/sso/sso.go:281` 的 `Saml2Acs`。

建议：

- 升级 `github.com/russellhaering/goxmldsig` 至 `v1.6.0` 或升级 `github.com/crewjam/saml` 到包含修复依赖的版本。
- 升级后重新运行 `GOTOOLCHAIN=go1.26.4 go run golang.org/x/vuln/cmd/govulncheck@latest ./...`。

### KBP-SEC-009 SSO 流程使用全局可变状态、固定 state，并信任 Referer 生成回调地址

严重级别：中

影响：

- SSO handler 使用包级全局变量保存 OIDC/SAML 配置，在并发登录或多用户登录时存在状态串扰风险。
- OAuth `state` 使用固定字符串，没有在回调中校验请求发起方状态，存在登录 CSRF 或账号绑定混淆风险。
- 回调地址由 `Referer` 字符串替换生成，部署在反向代理或存在可控 Referer 时，可能造成回调地址异常或安全边界不稳定。

证据：

- 全局状态：`internal/api/v1/sso/sso.go:20-21`
- 从 Referer 生成 callback/baseUrl：`internal/api/v1/sso/sso.go:103-120`
- 固定 `state`：`internal/api/v1/sso/sso.go:118`
- 回调仅取 `code` 并使用全局 `oauth2Config`，未校验 state：`internal/api/v1/sso/sso.go:144-148`

建议：

- 为每次登录生成随机 state/nonce，写入当前会话并在回调中校验。
- 不使用包级全局变量保存用户登录过程状态，改为会话级或服务端短期缓存。
- 回调地址从可信配置或服务端外部访问地址生成，不信任 `Referer`。

### KBP-SEC-010 多处关闭 TLS 证书校验

严重级别：中

影响：

- LDAP TLS、镜像仓库 HTTP client、webkubectl 生成 kubeconfig、非管理员 Kubernetes proxy transport 均关闭或跳过证书校验。
- 在中间人攻击或 DNS/网络被劫持时，可能泄露 LDAP 凭据、镜像仓库凭据、Kubernetes 用户证书或 token。

证据：

- LDAP TLS：`pkg/util/ldap/ldap_client.go:31-34`
- 镜像仓库：`pkg/util/imagerepo/repos/http_client.go:38-42`
- webkubectl kubeconfig：`internal/api/v1/webkubectl/webkubectl.go:65-75`
- proxy 非管理员 transport：`internal/api/v1/proxy/proxy.go:432-455`

建议：

- 默认启用证书校验，支持配置 CA bundle 或 per-cluster/per-repo CA。
- 仅允许管理员显式开启“不安全跳过校验”，并在 UI 和审计日志中明确标记。
- 对生成给用户的 kubeconfig 不应默认写入 `insecure-skip-tls-verify: true`。

### KBP-SEC-011 Cookie 会话缺少显式安全属性和 CSRF 防护

严重级别：中

影响：

- 服务端启用 cookie session，但 `sessions.Config` 只设置了 Cookie 名称、过期时间和 `AllowReclaim`，没有显式设置 `Secure`、`HttpOnly`、`SameSite`。
- 应用同时存在大量 cookie 可认证的 POST/PUT/DELETE 状态变更接口，未看到 CSRF token 或 Origin/Referer 校验。
- 如果部署在浏览器 cookie 登录模式下，存在跨站请求风险；若部署完全依赖 Authorization header，则风险降低。

证据：

- session 配置：`internal/server/server.go:169-171`
- 登录后设置 cookie：`internal/api/v1/session/session.go:173-180`

建议：

- 显式设置 `HttpOnly`、`Secure`、`SameSite=Lax/Strict`。
- 对状态变更接口增加 CSRF token 或严格 Origin 校验。
- 优先使用 Authorization header，并限制 cookie 认证的跨站使用范围。

### KBP-SEC-012 前端生产依赖存在多个 npm audit 告警

严重级别：中

影响：

- 多个前端包存在生产依赖漏洞告警。部分漏洞属于前端运行时或构建依赖，需要结合实际调用路径评估，但应纳入版本治理。
- `npm audit` 运行时环境提示 `NODE_TLS_REJECT_UNAUTHORIZED=0`，说明当前本机 npm 请求 TLS 校验被关闭；这属于本地环境风险，也建议修正。

扫描摘要：

- `web/dashboard`：25 个生产依赖告警，其中 critical 2、high 8。重点包括 `jsonpath-plus`、`form-data`、`axios`、`lodash`、`underscore`、`dompurify` 等。
- `web/kubepi`：14 个生产依赖告警，其中 critical 1、high 3。重点包括 `form-data`、`axios`、`js-yaml`、`nanotar` 等。
- `web/terminal`：8 个 high，集中在 Angular 16.2.x 相关包。
- `thirdparty/gotty/js`：生产依赖未发现 npm audit 告警。

证据：

- 直接依赖版本：`web/dashboard/package.json:18`、`web/dashboard/package.json:27-28`、`web/dashboard/package.json:32`、`web/dashboard/package.json:38`
- 直接依赖版本：`web/kubepi/package.json:17`、`web/kubepi/package.json:23-24`
- 直接依赖版本：`web/terminal/package.json:14-16`

建议：

- 优先升级 `axios`、`jsonpath-plus`、`dompurify`、`js-yaml`、`nanotar`、Angular 相关包。
- 升级后重新执行 `npm audit --omit=dev --json`，并配合前端构建和关键页面回归。
- 修正本机或 CI 环境中的 `NODE_TLS_REJECT_UNAUTHORIZED=0`。

## 已检查但未确认可利用的问题

- 前端 `v-html`：`web/dashboard/src/business/market-place/chart/detail/index.vue` 对 Markdown 渲染结果使用了 DOMPurify 清理，当前未确认存在直接 XSS；但 DOMPurify 自身依赖版本仍有 npm audit 告警，建议升级。
- `webkubectl/session` 下载 kubeconfig 为无认证路由，但 token 为一次性 UUID，创建 token 需要先认证；当前主要风险在 TLS 跳过校验和会话 token 生命周期，而不是直接未授权下载。
- Kubernetes proxy 路径对非管理员使用用户绑定证书，整体授权边界优于终端、日志、文件和 Helm 路径；但其非管理员 transport 仍关闭 TLS 校验，应按 `KBP-SEC-010` 加固。

## 验证与扫描记录

- 代码审计：按 API 路由、RBAC 白名单、默认角色、集群凭据使用点、Pod exec/log/file、SSO 登录、敏感模型返回逐项追踪。
- Go 依赖扫描：
  - 首次直接运行 `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` 因工具链切换到 Go 1.25 与项目 Go 1.26 不兼容而失败。
  - 强制 `GOTOOLCHAIN=go1.26.4` 后扫描成功，并确认 `GO-2026-4753` 可达。
  - `thirdparty/gotty` 使用同样方式扫描，未发现 Go 漏洞。
- npm 依赖扫描：
  - `npm audit --omit=dev --json` 已覆盖 `web/dashboard`、`web/kubepi`、`web/terminal`、`thirdparty/gotty/js`。

## 优先级建议

1. 立即移除 `charts`、`apps`、`pod` 白名单，并为这些资源补齐应用层 RBAC 和集群绑定校验。
2. 立即为集群、LDAP、用户、监控、镜像仓库等接口建立响应 DTO，默认不返回任何敏感字段。
3. 将终端、日志、Pod 文件、Helm 操作全部改为当前用户 Kubernetes 凭据执行，或限制为管理员专用。
4. 修复 SSO 自动创建用户固定密码问题，并迁移历史用户。
5. 升级 SAML 相关 Go 依赖，修复 `GO-2026-4753`。
6. 加固 TLS、cookie/CSRF、SSO state/nonce，并纳入自动化安全测试。
7. 分批升级前端依赖，先处理生产依赖中的 critical/high。
