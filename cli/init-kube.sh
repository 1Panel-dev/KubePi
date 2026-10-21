#!/bin/bash
set -e

arg1=$1
mkdir -p ~/.kube
code=`curl -w %{http_code} -s -o ~/.kube/config http://localhost/kubepi/api/v1/webkubectl/session?token=${arg1}`

if [[ $code -ne '200' ]];then
    echo "download kubeconfig failed"
    cat .kube/config
    exit
fi

current_context=`kubectl config current-context`
username=${current_context%@*}
cluster=${current_context#*@}

echo "Welcome to kubepi"
echo "Current cluster is ${username}"
echo "Current user is ${cluster}"

exec /bin/bash
