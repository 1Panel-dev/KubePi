# 使用 YAML 文件部署

## 快速开始

    sudo kubectl apply -f ./kubectl.yaml

## 获取访问地址

```sh
    # NodeIp
    export NODE_IP=$(kubectl get nodes -o jsonpath="{.items[0].status.addresses[0].address}")
    # NodePort
    export NODE_PORT=$(kubectl -n kube-system get services kubepi -o jsonpath="{.spec.ports[0].nodePort}")
    # Address
    echo http://$NODE_IP:$NODE_PORT
```

- 环境地址：http://$NODE_IP:$NODE_PORT
- 用户名：admin
- 密码：kubepi