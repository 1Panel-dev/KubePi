# AGENTS.md

本文件面向 AI 编码代理，介绍 KubePi 项目的架构、构建、测试与开发约定。阅读本文件前不需要任何项目背景知识。

## 项目概览

**KubePi**（[github.com/1Panel-dev/KubePi](https://github.com/1Panel-dev/KubePi)）是飞致云（FIT2CLOUD）出品的现代化 Kubernetes 多集群管理面板。管理员可导入多个 K8s 集群，通过 RBAC 权限控制把不同 cluster / namespace 的权限分配给指定用户；开发人员可用它管理集群中运行的应用并排查故障。

- 许可证：GPLv3
- 默认账号：`admin` / `kubepi`
- Go module 名：`github.com/KubeOperator/kubepi`（注意：仓库已迁移到 1Panel-dev 组织，但 module 名仍是 `KubeOperator`）
- 单进程交付：后端、前端静态资源、Web 终端工具最终打进一个二进制 `kubepi-server`（通过 `go:embed`），Docker 镜像 `1panel/kubepi` 暴露 80 端口

## 技术栈

### 后端（Go 1.22）

- Web 框架：[Iris v12](https://github.com/kataras/iris)（`github.com/kataras/iris/v12 v12.2.1`）
- CLI / 配置：cobra + viper，配置文件为 YAML（`conf/app.yml`，`kind: AppConfig`）
- 存储：**嵌入式 BoltDB**（经 `github.com/asdine/storm/v3` 封装），默认路径 `/var/lib/kubepi/db/kubepi.db`；无外置数据库
- Kubernetes：`k8s.io/client-go` v0.29.0、`k8s.io/kubectl`、`helm.sh/helm/v3` v3.14.2
- 认证：Iris session（cookie `SESS_COOKIE_KUBEPI`）+ JWT（HS256，`kataras/iris/v12/middleware/jwt`）；支持 LDAP（go-ldap）、SSO（OIDC / SAML）、MFA（xlzd/gotp）
- 日志：logrus（k8s.io/klog 输出被显式屏蔽，见 `internal/server/server.go` 的 `setUpLogger`）
- API 文档：swaggo 生成，运行时挂在 `/kubepi/swagger/`
- Web 终端：`thirdparty/gotty` 是 vendored 的 gotty 分支（`go.mod` 用 `replace` 指向本地目录），编译为独立二进制，运行时由主进程拉起并反向代理（`/kubepi/webkubectl/` → `localhost:8080`）
- i18n：`pkg/i18n` 内置 zh-CN / en-US 翻译表，默认语言 zh-CN

### 前端（3 个独立工程，构建产物经 `go:embed` 打进二进制）

| 目录              | 框架                           | 用途                                    |
| ----------------- | ------------------------------ | --------------------------------------- |
| `web/kubepi`    | Vue 2.7 + element-ui + vue-cli | 主管理界面（集群/用户/角色/镜像仓库等） |
| `web/dashboard` | Vue 2.7 + element-ui + echarts | 集群内 Kubernetes Dashboard             |
| `web/terminal`  | Angular 16 + xterm + sockjs    | Pod 终端                                |

三个前端各自 `npm run build` 后输出到 `cmd/server/web/{kubepi,dashboard,terminal}`（该目录已 gitignore，`vue.config.js` 里配置了 `outputDir`），`cmd/server/main.go` 用 `go:embed` 嵌入。前端开发服务器端口约定：kubepi 4300、terminal 4200、dashboard 4400；kubepi 与 terminal 的 dev server 把 `/kubepi/api` 代理到后端 2019，dashboard 代理到后端 80（见各自 `vue.config.js` / `proxy.config.json`）。

## 目录结构

```
cmd/server/            # 程序入口：main.go（cobra 命令 kubepi-server）、swagger docs、
                       #   script/{darwin,linux}/init-kube.sh 脚本、helper/ip/qqwry.dat IP 地址库
cli/                   # kubepi-cli：Node/TypeScript npm 包（bin: kubepi），
                       #   面向 AI agent 的 KubePi 命令行封装；含 skills/kubepi-cli/SKILL.md
conf/                  # app.yml / app_example.yml 配置样例
internal/              # 后端业务代码（Go internal 约定）
  api/v1/              #   HTTP 层：每个资源一个包（cluster、user、role、session、ldap、sso、
                       #   mfa、imagerepo、chart、file、proxy、ws、webkubectl、commons、system）
  model/v1/            #   数据模型（storm 实体，K8s 风格的 BaseModel/Metadata）
  service/v1/          #   业务逻辑层：接口 + 私有实现（NewService() 工厂）
  route/               #   路由注册（route.go → api/v1 的 AddV1Route）
  server/              #   Iris 应用装配：静态资源、session、结果/错误处理中间件、迁移
  config/              #   viper 读取 app.yml
migrate/               # BoltDB 数据迁移框架（版本号 + Handler，启动时按版本顺序执行）
pkg/                   # 可复用工具包：kubernetes 客户端封装、storm 查询扩展、i18n、
                       #   ldap/mfa/sso/helm/imagerepo/podtool 工具、IP 库、日志
thirdparty/gotty/      # vendored gotty（Web 终端后端）
web/                   # 3 个前端工程（见上表）
docs/deploy/           # 部署清单：compose、kubectl（charts 目录暂为空占位）
Dockerfile             # 三阶段构建（前端 → Go 二进制 → alpine 运行时）
Makefile               # 构建入口
```

### 后端分层约定（重要）
典型的请求链路是 `internal/api/v1/<resource>`（Handler，Iris）→ `internal/service/v1/<resource>`（Service 接口 + `NewService()` 返回私有 `service` 实现）→ storm DB。要点：

- Service 层统一通过 `common.DBOptions` 传递事务节点；不传时 `DefaultDBService.GetDB` 回退到全局 `server.DB()`
- Handler 不直接写响应体：把结果放进 `ctx.Values().Set("data", ...)`，由 `internal/server/server.go` 的统一结果中间件包装成 `{"success": true, "data": ...}`；错误用 `ctx.Values().Set("message", ...)` + 状态码，由统一错误中间件做 i18n 翻译后输出
- 路由注册在各资源包的 `Install(party)` 函数中，由 `internal/api/v1/v1.go` 的 `AddV1Route` 汇总；全局中间件链：JWT/session 认证 → 资源提取 → 角色加载 → RBAC 鉴权 → 资源名校验（`resourceNameInvalidHandler`）→ 操作日志
- API 路径前缀为 `/kubepi/api/v1/`；RBAC 的"资源名"取自 URL 第 4 段（`resourceExtractHandler`），`resourceWhiteList`（sessions、proxy、ws、charts、webkubectl、apps、mfa、pod）中的资源不做鉴权和操作日志
- 分页约定：query 参数 `pageNum` / `pageSize`（`pageHandler` 注入），条件查询统一走 `POST /<resource>/search`

## 构建与运行

### 常用命令（Makefile）

```bash
make build_web        # 构建全部 3 个前端（各自 npm install && npm run build）
make build_web_kubepi # 也可单独构建某一个前端（_dashboard / _terminal 同理）
make build_gotty      # 编译 thirdparty/gotty → dist/usr/local/bin/gotty
make build_bin        # 编译后端 → dist/usr/local/bin/kubepi-server（注意：依赖前端产物
                      #   已存在于 cmd/server/web/，否则 go:embed 编译失败）
make build_all        # build_web + build_gotty + build_bin
make build_docker     # 构建完整镜像
```

注意构建顺序：`go:embed` 要求 `cmd/server/web/kubepi` 等目录在编译前已生成，所以本地完整构建必须先 `make build_web` 再 `make build_bin`。

### 本地运行

```bash
go run cmd/server/main.go -c <配置目录> --server-bind-host 0.0.0.0 --server-bind-port 2019
```

- `-c/--config-path` 指向包含 `app.yml` 的目录；`conf/app_example.yml` 给出 server.bind、db.path、session.expires 三段样例，完整字段（含 jwt.key、logger.level）见配置模型 `internal/model/v1/config/config.go`
- 前端联调时在各 `web/*` 目录跑 `npm run serve`（terminal 用 `npm start`），dev server 已配置 `/kubepi/api`、`/webkubectl`、`/terminal` 的代理

### Docker / 部署

- `Dockerfile` 为三阶段构建：node:18.10 构建前端 → golang:1.22 编译 gotty 与 kubepi-server → alpine:3.16 运行时（内置 kubectl、k9s、kubectx/kubens、helm、fzf 等运维工具）
- 发布由 GitHub Actions 完成：`.github/workflows/release.yml`（打 `v*` tag 推送 multi-arch 镜像到 Docker Hub `1panel/kubepi`）、`pre-release.yml`（手动触发 RC）
- 部署清单在 `docs/deploy/`（compose / kubectl；charts 目录暂为空占位）
- 官方开发手册：[https://kubeoperator.io/docs/kubepi/dev_manual/](https://kubeoperator.io/docs/kubepi/dev_manual/)

## 测试

- 后端用标准 `go test`，测试文件以 `*_test.go` 放在被测包内。现有测试很少（如 `pkg/kubernetes`、`pkg/network/ip`、`pkg/util/ldap`、`pkg/util/podtool`、`internal/service/v1/rolebinding`），多为对工具函数的轻量测试；**项目没有统一的测试基建设施**（无 CI 测试阶段、无 mock 框架依赖、无测试数据库——存储本就是嵌入式 BoltDB，测试可直接开临时文件库）
- 运行：`go test ./...`（注意部分测试依赖 K8s 集群 / LDAP 等外部环境，失败不一定是代码问题）
- 前端：`web/kubepi`、`web/dashboard` 无测试脚本；`web/terminal` 有 `npm test`（karma + jasmine）

## 代码风格

- Go：标准 `gofmt` 风格；包注释与导出符号注释混用中英文，新代码跟随所在文件既有风格即可
- 命名约定：跨包的模型类型加包别名前缀导入（如 `v1User`、`v1RoleService`）；Handler 方法返回 `iris.Handler` 闭包
- 前端：vue-cli 自带 ESLint（`plugin:vue/essential` + `eslint:recommended`），`npm run lint` 检查
- Swagger 注解写在 Handler 方法上（`@Tags`、`@Router` 等）；重新生成文档的命令见 `cmd/server/main.go` 注释：`swag init -g "cmd/server/main.go" -o "cmd/server/docs" --parseDependency --parseInternal --parseDepth 2`

## 数据库迁移

- 迁移定义在 `migrate/v1/migrations.go`，每个迁移是 `{Version int, Message string, Handler func(tx)}`
- 启动时 `migrate.RunMigrate` 读取 DB 中 `current_db_version`，在单个事务中按版本升序执行所有更新版本的迁移，失败即回滚并退出进程
- 新增迁移只需往 `Migrations` 列表追加更大版本号的条目，**不要修改已发布的迁移**

## 安全注意事项

- 鉴权有两条路径：session cookie 和 `Authorization: Bearer <JWT>`，JWT 密钥来自配置 `spec.jwt.key`；改动认证逻辑时两条路径都要顾及（见 `internal/api/v1/v1.go` 的 `WarpedJwtHandler` / `authHandler`）
- RBAC 在 URL 段级别做资源匹配，新增路由时注意资源名会暴露为权限项（`apiResourceHandler` 自动枚举路由生成权限列表）
- 密码用 bcrypt 存储（`internal/service/v1/user`）；LDAP 绑定、SSO 回调、镜像仓库凭据等敏感配置均落库于本地 BoltDB
- 上报安全问题：support@fit2cloud.com（见 `SECURITY.md`）；提交 issue/PR 时脱敏（用户名、密码、IP 等用 `REDACTED` 代替）

## 贡献流程

- PR 拆小、渐进式提交，每个合并到 master 的变更都不能破坏主干（见 `CONTRIBUTING.md`）
- Issue / PR 模板在 `.github/` 下
