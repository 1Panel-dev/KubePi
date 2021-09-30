<p align="center"><a href="https://kubeoperator.io"><img src="https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/img/logo-red.png" alt="kubepi" width="300" /></a></p>
<h3 align="center">简单易用的开源 Kubernetes 可视化管理面板</h3>
<p align="center">
  <a href="http://www.apache.org/licenses/LICENSE-2.0"><img src="https://img.shields.io/github/license/kubeoperator/kubepi?color=%231890FF&style=flat-square" alt="License: Apache License v2"></a>
  <a href="https://app.codacy.com/gh/kubeoperator/kubepi?utm_source=github.com&utm_medium=referral&utm_content=kubeoperator/kubepi&utm_campaign=Badge_Grade_Dashboard"><img src="https://app.codacy.com/project/badge/Grade/da67574fd82b473992781d1386b937ef" alt="Codacy"></a>
  <a href="https://hub.docker.com/r/kubeoperator/kubepi-server"><img src="https://img.shields.io/docker/pulls/kubeoperator/kubepi-server" alt="Docker Pulls"></a>
  <a href="https://github.com/KubeOperator/KubePi"><img src="https://img.shields.io/github/stars/KubeOperator/KubePi" alt="Stars"></a>
</p>
<hr />

KubePi 是一款简单易用的开源 Kubernetes 可视化管理面板。

KubePi 允许管理员导入多个 Kubernetes 集群，并且通过权限控制，将不同 cluster、namespace 的权限分配给指定用户。它允许开发人员管理 Kubernetes 集群中运行的应用程序并对其进行故障排查，供开发人员更好地处理 Kubernetes 集群中的复杂性。

### UI 展示

![UI展示](https://kubeoperator.oss-cn-beijing.aliyuncs.com/kubepi/img/kubepi-demo.gif)

### 快速开始

    sudo docker run -d --restart=unless-stopped -p 80:80 kubeoperator/kubepi-server

打开浏览器访问http://localhost

### 微信交流群

![wechat-group](https://kubeoperator.io/docs/img/wechat-group.png)

### 致谢

- [Vue](https://cn.vuejs.org) 前端框架
- [FIT2CLOUD UI](https://github.com/fit2cloud-ui/fit2cloud-ui/) FIT2CLOUD UI组件库
- [Vue-element-admin](https://github.com/PanJiaChen/vue-element-admin) 项目脚手架

### License & Copyright

Copyright (c) 2014-2021 FIT2CLOUD 飞致云

[https://www.fit2cloud.com](https://www.fit2cloud.com)<br>

KubePi is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
