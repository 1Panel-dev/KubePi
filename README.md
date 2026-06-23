<p align="center"><a href="https://github.com/1Panel-dev/KubePi"><img src="https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/img/logo-red.png" alt="kubepi" width="300" /></a></p>
<P align="center"><b>KubePi</b> [kubəˈpaɪ]，一个现代化的 K8s 面板。</P>
<p align="center">
  <a href="https://github.com/1Panel-dev/KubePi/releases"><img src="https://img.shields.io/github/v/release/1Panel-dev/KubePi" alt="GitHub release"></a>
  <a href="https://github.com/1Panel-dev/KubePi"><img src="https://img.shields.io/github/stars/1Panel-dev/KubePi?color=%231890FF&style=flat-square" alt="GitHub Stars"></a>
  <a href="https://hub.docker.com/r/1panel/kubepi"><img src="https://img.shields.io/docker/pulls/1panel/kubepi?label=downloads" alt="Docker Pulls"></a>
</p>
<hr />

## KubePi 是什么？

KubePi 是一个现代化的 K8s 面板。KubePi 允许管理员导入多个 Kubernetes 集群，并且通过权限控制，将不同 cluster、namespace 的权限分配给指定用户；允许开发人员管理 Kubernetes 集群中运行的应用程序并对其进行故障排查，供开发人员更好地处理 Kubernetes 集群中的复杂性。

## 快速开始

```
docker run --privileged -d --restart=unless-stopped -p 80:80 1panel/kubepi

# 用户名: admin
# 密码: kubepi
```

你也可以通过 [1Panel 应用商店](https://apps.fit2cloud.com/1panel) 快速部署 KubePi。

使用手册请参考: [https://github.com/1Panel-dev/KubePi/wiki](https://github.com/1Panel-dev/KubePi/wiki)。

## 版本兼容性

KubePi v2.0.0 使用 Go 1.26 和 Kubernetes client-go v0.36 构建。

- 最低支持 Kubernetes 版本：v1.24。
- 推荐 Kubernetes 版本：v1.34 - v1.36。
- Kubernetes v1.23 及更早版本不再作为 v2.0.0 的支持范围。
- PodSecurityPolicy 已在 Kubernetes v1.25 移除，KubePi 仅在目标集群仍提供 `policy/v1beta1` API 时展示相关入口。

## UI 展示

![UI展示](https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/img/02-dashboard.png)

## 飞致云的其他明星项目

- [1Panel](https://github.com/1panel-dev/1panel/) - 现代化、开源的 Linux 服务器运维管理面板
- [Cordys CRM](https://github.com/1Panel-dev/CordysCRM) - 新一代的开源 AI CRM 系统
- [JumpServer](https://github.com/jumpserver/jumpserver/) - 广受欢迎的开源堡垒机
- [MaxKB](https://github.com/1Panel-dev/MaxKB/) - 基于 LLM 大语言模型的开源知识库问答系统
- [DataEase](https://github.com/dataease/dataease/) - 人人可用的开源数据可视化分析工具
- [MeterSphere](https://github.com/metersphere/metersphere/) - 开源持续测试工具
- [Halo](https://github.com/halo-dev/halo/) - 强大易用的开源建站工具

## License

本仓库遵循 [FIT2CLOUD Open Source License](LICENSE) 开源协议，该许可证本质上是 GPLv3，但有一些额外的限制。

你可以基于 KubePi 的源代码进行二次开发，但是需要遵守以下规定：

- 不能替换和修改 KubePi 的 Logo 和版权信息；
- 二次开发后的衍生作品必须遵守 GPL V3 的开源义务。

如需商业授权，请联系：`support@fit2cloud.com`。
