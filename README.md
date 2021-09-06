# KubePi

KubePi 是一款开源的容器管理平台。它允许用户管理 Kubernetes 集群中运行的应用程序并对其进行故障排查，供开发人员更好地了解 Kubernetes 集群的复杂性。

## KubePi 的功能：

- 多集群管理：集群管理员可以导入多个 Kubernetes 集群，并且通过权限控制，将不同 cluster、namespace 的权限分配给指定用户
- 工作负载可视化：提供了图形化的工作负载编辑界面，用户可轻松完成针对容器的编排任务
- 日志、终端：支持查看 pod 日志和针对 pod 中的容器执行命令，以便进行故障排除或监控

## 快速开始

    sudo docker run -v ./kubepi/:/var/lib/kubepi -p 2019:2019 kubeoperator/kubepi-server:v1.0

- 环境地址：http://x.x.x.x:2019
- 用户名：admin
- 密码：admin123

## 其他安装方式

- [kubernetes](docs/deploy/kubectl)
- [docker-compose](docs/deploy/compose)

## 微信交流群

![wechat-group](https://kubeoperator.io/docs/img/wechat-group.png)

## 版本说明

KubePi 版本号命名规则为：v大版本.功能版本.Bug修复版本。比如：

```
v1.0.1 是 v1.0.0 之后的Bug修复版本；
v1.1.0 是 v1.0.0 之后的功能版本。
```

像其它优秀开源项目一样，KubePi 将每月发布一个功能版本。

## 致谢

- [Vue](https://cn.vuejs.org) 前端框架
- [FIT2CLOUD UI](https://github.com/fit2cloud-ui/fit2cloud-ui/) FIT2CLOUD UI组件库
- [Vue-element-admin](https://github.com/PanJiaChen/vue-element-admin) 项目脚手架

## License & Copyright

Copyright (c) 2014-2021 FIT2CLOUD 飞致云

[https://www.fit2cloud.com](https://www.fit2cloud.com)<br>

KubePi is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
