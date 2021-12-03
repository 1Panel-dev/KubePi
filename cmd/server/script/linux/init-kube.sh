#!/bin/bash
set -e

echo "export TERM=xterm-256color" >> /root/.bashrc
echo "source /usr/share/bash-completion/bash_completion" >> /root/.bashrc
echo 'source <(kubectl completion bash)' >> /root/.bashrc
echo 'complete -F __start_kubectl k' >> /root/.bashrc

if [ "${WELCOME_BANNER}" ]; then
    echo ${WELCOME_BANNER}
fi

arg1=$1


mkdir -p /nonexistent
mount -t tmpfs -o size=10m tmpfs /nonexistent
cd /nonexistent
cp /root/.bashrc ./
cp /etc/vim/vimrc.local .vimrc
echo 'source /opt/kubectl-aliases/.kubectl_aliases' >> .bashrc
mkdir -p .kube

export HOME=/nonexistent


code=`curl -w %{http_code} -s -o ~/.kube/config http://localhost/kubepi/api/v1/webkubectl/session?token=${arg1}`

if [[ $code -ne '200' ]];then
    echo "download kubeconfig failed"
    cat .kube/config
    exit
fi

current_context=`kubectl config current-context`
cluster=${current_context%@*}
username=${current_context#*@}

echo -e PS1=\"'['${username}@${cluster}']$ '\"  >> .bashrc

echo "Welcome to kubepi web terminal, try kubectl --help."



chown -R nobody:nogroup .kube
export TMPDIR=/nonexistent

envs=`env`
for env in ${envs[@]}; do
    if [[ $env == GOTTY* ]];
    then
        unset ${env%%=*}
    fi
done

unset WELCOME_BANNER PPROF_ENABLED KUBECTL_INSECURE_SKIP_TLS_VERIFY SESSION_STORAGE_SIZE KUBECTL_VERSION

exec su -s /bin/bash nobody
exec /bin/bash
