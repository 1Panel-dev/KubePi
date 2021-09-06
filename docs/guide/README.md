# 导入集群

## 获取 apiserver 地址:

    cat .kube/config | grep server: | awk '{print $2}'

> 注意: 如果 server IP 为 127.0.0.1，需要将 IP 替换为任意 master 节点 IP

## 获取 token:

    echo $(kubectl -n kube-system get secret $(kubectl -n kube-system get secret | grep kubepi-user | awk '{print $1}') -o go-template='{{.data.token}}' | base64 -d)
