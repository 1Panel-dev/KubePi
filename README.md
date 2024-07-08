<p align="center"><a href="https://kubeoperator.io"><img src="https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/img/logo-red.png" alt="kubepi" width="300" /></a></p>
<P align="center"><b>KubePi</b> [kubəˈpaɪ]，一个现代化的 K8s 面板。</P>
<p align="center">
  <a href="http://www.apache.org/licenses/LICENSE-2.0"><img src="https://img.shields.io/github/license/kubeoperator/kubepi?color=%231890FF&style=flat-square" alt="License: Apache License v2"></a>
  <a href="https://app.codacy.com/gh/kubeoperator/kubepi?utm_source=github.com&utm_medium=referral&utm_content=kubeoperator/kubepi&utm_campaign=Badge_Grade_Dashboard"><img src="https://app.codacy.com/project/badge/Grade/da67574fd82b473992781d1386b937ef" alt="Codacy"></a>
  <a href="https://hub.docker.com/r/kubeoperator/kubepi-server"><img src="https://img.shields.io/docker/pulls/kubeoperator/kubepi-server" alt="Docker Pulls"></a>
  <a href="https://github.com/1Panel-dev/KubePi"><img src="https://img.shields.io/github/stars/1Panel-dev/KubePi?color=%231890FF&style=flat-square" alt="GitHub Stars"></a>
</p>
<hr />

KubePi 是一个现代化的 K8s 面板。

KubePi 允许管理员导入多个 Kubernetes 集群，并且通过权限控制，将不同 cluster、namespace 的权限分配给指定用户。它允许开发人员管理 Kubernetes 集群中运行的应用程序并对其进行故障排查，供开发人员更好地处理 Kubernetes 集群中的复杂性。

### UI 展示

![UI展示](https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/img/02-dashboard.png)

### 快速开始

```
docker run --privileged -d --restart=unless-stopped -p 80:80 1panel/kubepi

# 用户名: admin
# 密码: kubepi
```

你也可以通过 [1Panel 应用商店](https://apps.fit2cloud.com/1panel) 快速部署 KubePi。

## 飞致云的其他明星项目

- [1Panel](https://github.com/1panel-dev/1panel/) - 现代化、开源的 Linux 服务器运维管理面板
- [JumpServer](https://github.com/jumpserver/jumpserver/) - 广受欢迎的开源堡垒机
- [DataEase](https://github.com/dataease/dataease/) - 人人可用的开源数据可视化分析工具
- [MeterSphere](https://github.com/metersphere/metersphere/) - 开源自动化测试平台
- [Halo](https://github.com/halo-dev/halo/) - 强大易用的开源建站工具
- [MaxKB](https://github.com/1Panel-dev/MaxKB/) - 基于 LLM 大语言模型的开源知识库问答系统

### License & Copyright

Copyright (c) 2014-2024 FIT2CLOUD 飞致云

[https://www.fit2cloud.com](https://www.fit2cloud.com)<br>

KubePi is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
