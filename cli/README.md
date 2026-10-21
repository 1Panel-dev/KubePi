# kubepi-cli

KubePi 命令行工具 —— 在终端里访问和操作 [KubePi](https://github.com/1Panel-dev/KubePi)（Kubernetes 多集群管理面板），对 AI Agents 友好。

## 特性

- 覆盖 KubePi 全部管理资源：集群、成员、角色、用户、镜像仓库、Chart、LDAP、SSO、监控、日志
- `k8s` 透传命令：经 KubePi 代理操作集群内任意 Kubernetes API
- 登录凭证缓存，JWT 过期自动重登（KubePi JWT 默认 10 分钟）
- 默认 Markdown 表格输出，`--format=json` 输出结构化 JSON
- `--pick` / `--filter` / `--search` / `--sort` / `--page` / `--all` / `--limit` 数据处理
- 结构化错误码（E1001/E1004/...），便于 AI 识别与恢复
- 随包附带 [Agent Skill](skills/kubepi-cli/SKILL.md)

## 安装

```bash
npm install -g kubepi-cli        # 安装最新版
npm install -g kubepi-cli@0.1.0  # 指定版本
kubepi help
```

要求 Node.js >= 18。源码开发方式：`npm install && npm run build && npm link`，不 link 也可以直接用 `node bin/kubepi.js`。

## 快速开始

```bash
# 登录（凭证缓存到 ~/.config/kubepi/config.json）
kubepi login -s http://127.0.0.1:2019 -u admin -p kubepi

# 列出集群 / 用户
kubepi cluster
kubepi user --search=tom --format=json

# 导入集群
kubepi cluster create --name=dev --mode=configFile --kubeconfig=~/.kube/config

# 操作集群内 K8s 资源
kubepi k8s dev /api/v1/pods --namespace=default
```

## 开发

```bash
bun install
bun run dev -- help     # 直接跑 TS 源码
bun test                # 单元测试
npm run typecheck       # 类型检查
npm run build           # 产出 dist/index.js
```

兼容 KubePi v1.9 与 v2.x；monitor、user reset-password 等少数功能需要 KubePi >= v2.0。

## 许可证

GPL-3.0（与 KubePi 主项目一致）
