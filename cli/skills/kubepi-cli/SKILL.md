---
name: kubepi-cli
description: 通过 kubepi 命令行工具查询和操作 KubePi（Kubernetes 多集群管理面板），覆盖集群、集群成员、集群角色、用户、角色、镜像仓库、Chart 仓库与应用、LDAP、SSO、监控、日志等模块的增删改查，以及通过 k8s 透传命令操作集群内任意 Kubernetes 资源。当用户提到 KubePi、kubepi、导入集群、管理集群成员、给用户授权、查 Pod/Deployment 等 K8s 资源时使用本技能。
license: GPL-3.0
metadata:
  author: 1Panel-dev
  repository: https://github.com/1Panel-dev/KubePi.git
  keywords: [kubepi, kubernetes, cli, multi-cluster]
  version: 0.1.0
---
# KubePi CLI

通过 `kubepi` 命令行工具查询和操作 KubePi。CLI 自动处理认证（JWT 过期自动重登）、分页，支持过滤/搜索/排序/字段摘取。

## 前置准备

### 安装

```bash
npm install -g kubepi-cli

# 安装技能到指定 Agent
kubepi add-skill kimi-code

# 安装到全部支持的 Agent
kubepi add-skill all

# 或导出到自定义目录
kubepi add-skill --output ./skills
```

支持的 Agent：`claude-code`、`cursor`、`cherry-studio`、`codex`、`opencode`、`vscode`、`antigravity`、`gemini`、`kimi-code`、`zcode`。

源码开发方式：`cd cli && npm install && npm run build && npm link`，或 `bun install && bun run build && bun link`。构建产物为 `dist/index.js`，`bin/kubepi.js` 为入口。也可以不 link，直接 `node bin/kubepi.js` 使用。

### 认证

首次执行：

```bash
kubepi login -s http://127.0.0.1:2019 -u admin -p kubepi
```

环境变量（优先级低于命令行参数）：

| 变量                | 说明                                     |
| ------------------- | ---------------------------------------- |
| `KUBEPI_URL`      | KubePi 服务地址                          |
| `KUBEPI_USERNAME` | 用户账号                                 |
| `KUBEPI_PASSWORD` | 密码                                     |
| `KUBEPI_TOKEN`    | 直接指定 JWT Token（有此变量可省略密码） |

登录成功后凭证缓存在 `~/.config/kubepi/config.json`（0600 权限）。KubePi 签发的 JWT 默认 10 分钟过期，CLI 会在过期前用缓存账密自动重新登录，无需人工干预。

### 凭证安全

- 用户尚未登录时，不要在对话里收集账号密码。让用户直接在终端执行 `kubepi login`，由用户自行输入凭证。
- 严禁读取本地凭证：`KUBEPI_PASSWORD` / `KUBEPI_TOKEN` 环境变量、`~/.config/kubepi/config.json` 配置文件。所有 KubePi 数据均通过 `kubepi` 命令获取，凭证由 CLI 内部处理。

## 命令格式

使用简写方式（推荐）：

| 操作 | 命令                                            |
| ---- | ----------------------------------------------- |
| 列表 | `kubepi <module>`                             |
| 详情 | `kubepi <module> <name>`                      |
| 创建 | `kubepi <module> create --field=value`        |
| 更新 | `kubepi <module> update <name> --field=value` |
| 删除 | `kubepi <module> delete <name>`               |
| 动作 | `kubepi <module> <action> [args]`             |
| 帮助 | `kubepi <module> help`                        |

也支持 `--data='JSON'` 传入完整 JSON 数据；`--field` 支持点路径（如 `--spec.connect.forward.apiServer=x`）。

## 模块与操作速查

| 模块名        | 中文                   | 支持的操作                                                                 |
| ------------- | ---------------------- | -------------------------------------------------------------------------- |
| cluster       | 集群                   | CRUD + search + namespaces / members / repos 动作                          |
| clustermember | 集群成员               | CRUD（需 `--cluster=<name>`）                                            |
| clusterrole   | 集群内 K8s ClusterRole | CUD + 列表（需 `--cluster=<name>`）                                      |
| clusterrepo   | 集群绑定的 Chart 仓库  | CD + 列表 + detail（需 `--cluster=<name>`）                              |
| user          | 用户                   | CRUD + reset-password（≥ v2.0）                                           |
| role          | 角色                   | CRUD + search                                                              |
| imagerepo     | 镜像仓库               | CRUD + images / search-images / repositories                               |
| chartrepo     | Chart 仓库             | CRUD + sync（需 `--cluster=<name>`）                                     |
| chart         | Chart                  | 列表 / detail / install / installed / uninstall（需 `--cluster=<name>`） |
| ldap          | LDAP 配置              | 列表 + create / update / sync / test-connect / test-login / import         |
| sso           | SSO 配置               | get / set / status                                                         |
| monitor       | 监控（≥ v2.0）        | grafana* + metrics*                                                        |
| log           | 日志                   | login / operation                                                          |
| k8s           | 通用透传               | `kubepi k8s <cluster> <apiPath>` 任意 K8s API 读写                       |

> name 即 ID：KubePi 资源以 name 为唯一标识，没有数字 ID。
> 标注 ≥ v2.0 的功能在老版本服务端返回 404（E2002）。

## AI 使用策略

### 输出格式

- 展示给用户：不加 `--format` 参数，默认输出 Markdown 表格（列表）或 k-v 列表（单个对象）
- 需要程序化处理：加 `--format=json`，返回结构化 JSON（已解包 `data` 字段）

### 交互确认

AI 场景下执行删除操作时加 `--yes` 跳过确认提示：

```bash
kubepi user delete tom --yes
```

非 TTY 环境下删除不带 `--yes` 会直接报 E2003。

### 全局选项位置无关

`--format`、`--pick`、`--filter`、`--search`、`--sort`、`--page`、`--page-size`、`--all`、`--limit`、`--yes`、`--data`、`--server`、`--username`、`--password`、`--token` 可写在命令行任意位置。

### 不知道名字时

先查列表获取 name，再操作具体对象：

```bash
kubepi cluster --pick=name,status.phase
kubepi user admin
```

### 写操作前确认

执行创建、更新、删除等写操作前，先向用户确认操作内容。用户明确要求不确认时可跳过。

### 更新操作自动补全

执行 `update` 时，CLI 会先 GET 当前对象，把未显式传入的字段用现值填充后再 PUT。只需传想改的字段即可。

## 数据处理

### 摘取字段（支持点路径）

```bash
kubepi cluster --pick=name,status.phase,spec.connect.forward.apiServer
```

### 过滤

```bash
kubepi user --filter='type:LOCAL'                      # : 等于
kubepi user --filter='name~adm'                        # ~ 包含
kubepi user --filter='isAdmin:false,type:LOCAL'        # 逗号 = AND
kubepi user --filter='name:admin' --filter='name:tom'  # 多个 --filter = OR
```

支持的运算符：`:` `=` 等于、`!=` 不等于、`>` `<` `>=` `<=`、`~` 包含、`!~` 不包含。可服务端的条件会下推到 `/search` 接口，其余在客户端精确过滤。

### 模糊搜索

```bash
kubepi user --search=admin
kubepi chart --cluster=dev --search=nginx
```

### 排序

```bash
kubepi user --sort=createAt_desc
```

### 分页

```bash
kubepi user --page=1 --page-size=50
kubepi user --all            # 拉取全部页
kubepi user --limit=10       # 只取前 10 条
```

## 常用操作示例

### 导入集群并授权

```bash
kubepi cluster create --name=dev --mode=bearer --api-server=https://1.2.3.4:6443 --bearer-token=<token>
kubepi cluster create --name=prod --mode=configFile --kubeconfig=~/.kube/config
kubepi clustermember create --cluster=dev --name=tom --cluster-role=cluster-admin
kubepi clustermember create --cluster=dev --name=tom --namespace-role=default:view,edit
```

### 创建角色并授权给用户

```bash
kubepi role create --name=dev --rules='[{"resource":["clusters","proxy"],"resourceNames":["dev"],"verbs":["get","list"]}]'
kubepi user create --name=tom --email=tom@example.com --password=Secret123 --role=dev
kubepi user update tom --role=ops --role=dev          # 覆盖式替换 roles
kubepi user update tom --password=NewSecret --old-password=Secret123
```

### 配置 Chart 仓库并安装应用

```bash
kubepi chartrepo create --cluster=dev --name=bitnami --url=https://charts.bitnami.com/bitnami
kubepi chartrepo sync bitnami --cluster=dev
kubepi chart --cluster=dev --repo=bitnami --search=nginx
kubepi chart detail nginx --cluster=dev --repo=bitnami --version=15.0.0
kubepi chart install --cluster=dev --name=my-nginx --chart-name=nginx --repo=bitnami --chart-version=15.0.0 --namespace=default
kubepi chart installed --cluster=dev
```

### 操作集群内 K8s 资源（k8s 透传）

```bash
kubepi k8s dev /api/v1/pods --namespace=default
kubepi k8s dev /api/v1/pods --search=nginx --page=1 --page-size=20
kubepi k8s dev /apis/apps/v1/deployments --namespace=default --format=json
kubepi k8s dev /api/v1/namespaces --method=POST --data='{"metadata":{"name":"demo"}}'
kubepi k8s dev /api/v1/namespaces/demo --method=DELETE --yes
```

`<path>` 是 Kubernetes API 路径（如 `/api/v1/pods`、`/apis/apps/v1/deployments`）。列表结果默认以表格展示 `metadata.name/namespace/creationTimestamp`；`--format=json` 输出服务端原始返回。

### 日志审计

```bash
kubepi log login --search=admin
kubepi log operation --filter='operator:admin' --page=1 --page-size=20
```

### 查看帮助

```bash
kubepi help                 # 查看所有模块
kubepi cluster help         # 查看集群模块的参数和操作
kubepi chart install help   # 查看 install 动作的参数
```

## 意图识别

| 用户意图                      | CLI 命令                                                                                        |
| ----------------------------- | ----------------------------------------------------------------------------------------------- |
| 登录/切换账号                 | `kubepi login` / `kubepi logout` / `kubepi profile`                                       |
| 所有集群                      | `kubepi cluster`                                                                              |
| 导入/接入集群                 | `kubepi cluster create ...`                                                                   |
| 集群的 namespace              | `kubepi cluster namespaces <name>`                                                            |
| 集群成员授权                  | `kubepi clustermember --cluster=<c>` / `create --cluster=<c> --name=<u> --cluster-role=...` |
| 集群内 ClusterRole            | `kubepi clusterrole --cluster=<c>`                                                            |
| 创建用户                      | `kubepi user create ...`                                                                      |
| 重置密码                      | `kubepi user reset-password <name> --password=...`（≥ v2.0）                                 |
| 创建角色                      | `kubepi role create ...`                                                                      |
| 镜像仓库                      | `kubepi imagerepo ...`                                                                        |
| Chart 仓库/应用商店           | `kubepi chartrepo ...` / `kubepi chart ...`                                                 |
| 安装应用                      | `kubepi chart install ...`                                                                    |
| 查 Pod/Deployment 等 K8s 资源 | `kubepi k8s <cluster> /api/v1/...`                                                            |
| LDAP 配置                     | `kubepi ldap ...`                                                                             |
| SSO 配置                      | `kubepi sso get` / `kubepi sso set ...`                                                     |
| 监控/Grafana                  | `kubepi monitor ...`（≥ v2.0）                                                               |
| 登录/操作日志                 | `kubepi log login` / `kubepi log operation`                                                 |

## 错误处理

| 错误码 | 含义                          | 处理方式                                                           |
| ------ | ----------------------------- | ------------------------------------------------------------------ |
| E1001  | 未登录/凭证缺失               | 执行 `kubepi login`                                              |
| E1004  | 登录失败或 Token 失效         | 检查账密后执行 `kubepi login` 重新登录                           |
| E1005  | 账号启用了 MFA                | 换用未启用 MFA 的账号（KubePi ≥ v2.0 拒绝 MFA 账号的 JWT 登录）   |
| E2001  | 模块不存在                    | 执行 `kubepi help` 查看可用模块                                  |
| E2002  | 对象不存在                    | 检查 name；部分功能（monitor、reset-password）需要 KubePi ≥ v2.0  |
| E2003  | 缺少必要参数                  | 执行 `kubepi <module> help` 或 `kubepi <module> <action> help` |
| E2004  | 参数非法（格式错误/类型不对） | 检查参数取值，如数值型 flag 传了非数字、JSON 无法解析              |
| E2006  | 无权限                        | 检查当前账号的 RBAC 权限                                           |
| E5001  | 请求失败/超时                 | 检查网络或 KubePi 服务状态                                         |

## 注意事项

- 不确定模块参数时，先执行 `kubepi <module> help`；不确定动作参数时，先执行 `kubepi <module> <action> help`
- 兼容 KubePi v1.9 与 v2.x（端点路径一致）；≥ v2.0 的独有功能已在模块表中标注
- 集群认证模式：`bearer`（token）、`certificate`（客户端证书）、`configFile`（kubeconfig）
- `mfa` 接口基于浏览器 session，CLI 无法使用，请在 Web UI 操作
- 集群创建是异步初始化，`status.phase` 变为 `Completed` 后才可用
- `kubepi logout` 只清除本地缓存凭证；服务端 JWT 在其 10 分钟有效期内自然过期
- `kubepi profile` 的数据来自 JWT 本地解码（`GET /sessions` 是浏览器 session 专用接口）
- `login -p` 明文密码会进入 shell 历史；交互使用时建议省略 `-p`，走密码提示输入（不回显）
