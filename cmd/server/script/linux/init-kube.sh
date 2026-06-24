#!/bin/bash
set -u

start_shell() {
    export TMPDIR=/nonexistent
    if command -v su >/dev/null 2>&1 && id nobody >/dev/null 2>&1; then
        su -s /bin/bash nobody
        status=$?
        if [[ $status -eq 0 ]]; then
            exit 0
        fi
        echo "warning: can not start shell as nobody, fallback to bash"
    fi
    exec /bin/bash
}

fail_and_shell() {
    echo "$1"
    start_shell
}

echo "export TERM=xterm-256color" >> /root/.bashrc 2>/dev/null || true
echo "source /usr/share/bash-completion/bash_completion" >> /root/.bashrc 2>/dev/null || true
echo 'source <(kubectl completion bash)' >> /root/.bashrc 2>/dev/null || true
echo 'complete -F __start_kubectl k' >> /root/.bashrc 2>/dev/null || true

if [ "${WELCOME_BANNER:-}" ]; then
    echo ${WELCOME_BANNER}
fi

arg1=${1:-}

if [[ -z "${arg1}" ]]; then
    fail_and_shell "missing terminal session token"
fi

mkdir -p /nonexistent || fail_and_shell "can not create terminal home"
mount -t tmpfs -o size=10m tmpfs /nonexistent 2>/tmp/kubepi-mount.err || {
    echo "warning: can not mount tmpfs for terminal home"
    cat /tmp/kubepi-mount.err
}
cd /nonexistent || fail_and_shell "can not enter terminal home"
cp /root/.bashrc ./ 2>/dev/null || touch .bashrc
cp /etc/vim/vimrc.local .vimrc 2>/dev/null || true
echo 'source /opt/kubectl-aliases/.kubectl_aliases' >> .bashrc
mkdir -p .kube || fail_and_shell "can not create kubeconfig directory"

export HOME=/nonexistent

session_url=${KUBEPI_WEBKUBECTL_SESSION_URL:-http://127.0.0.1/kubepi/api/v1/webkubectl/session}
code=$(curl --connect-timeout 5 --max-time 20 -w "%{http_code}" -s -o ~/.kube/config "${session_url}?token=${arg1}")
curl_status=$?

if [[ $curl_status -ne 0 || "${code}" != "200" ]];then
    echo "download kubeconfig failed (curl: ${curl_status}, http: ${code:-000})"
    rm -f .kube/config
    start_shell
fi

current_context=$(kubectl config current-context 2>/tmp/kubepi-kubectl.err)
if [[ $? -ne 0 || -z "${current_context}" ]]; then
    echo "load kubeconfig failed"
    cat /tmp/kubepi-kubectl.err
    start_shell
fi
cluster=${current_context%@*}
username=${current_context#*@}

echo -e PS1=\"'['${username}@${cluster}']$ '\"  >> .bashrc

echo "Welcome to kubepi web terminal, try kubectl --help."



chown -R nobody:nogroup /nonexistent 2>/dev/null || chown -R nobody /nonexistent 2>/dev/null || true
export TMPDIR=/nonexistent

envs=`env`
for env in ${envs[@]}; do
    if [[ $env == GOTTY* ]];
    then
        unset ${env%%=*}
    fi
done

unset WELCOME_BANNER PPROF_ENABLED KUBECTL_INSECURE_SKIP_TLS_VERIFY SESSION_STORAGE_SIZE KUBECTL_VERSION

start_shell
