#!/bin/bash
set -u

start_shell() {
    export TMPDIR=/nonexistent
    if ! command -v su >/dev/null 2>&1 || ! id nobody >/dev/null 2>&1; then
        echo "can not start terminal as nobody" >&2
        exit 1
    fi
    exec su -s /bin/bash nobody
    echo "can not start terminal as nobody" >&2
    exit 1
}

fail_and_exit() {
    echo "$1" >&2
    exit 1
}

arg1=${1:-}

if [[ -z "${arg1}" ]]; then
    fail_and_exit "missing terminal session token"
fi

echo "export TERM=xterm-256color" >> /root/.bashrc 2>/dev/null || true
echo "source /usr/share/bash-completion/bash_completion" >> /root/.bashrc 2>/dev/null || true
echo 'source <(kubectl completion bash)' >> /root/.bashrc 2>/dev/null || true
echo 'complete -F __start_kubectl k' >> /root/.bashrc 2>/dev/null || true

if [ "${WELCOME_BANNER:-}" ]; then
    echo ${WELCOME_BANNER}
fi

mkdir -p /nonexistent || fail_and_exit "can not create terminal home"
mount -t tmpfs -o size=10m tmpfs /nonexistent 2>/tmp/kubepi-mount.err || {
    echo "warning: can not mount tmpfs for terminal home"
    cat /tmp/kubepi-mount.err
}
cd /nonexistent || fail_and_exit "can not enter terminal home"
cp /root/.bashrc ./ 2>/dev/null || touch .bashrc
cp /etc/vim/vimrc.local .vimrc 2>/dev/null || true
echo 'source /opt/kubectl-aliases/.kubectl_aliases' >> .bashrc
mkdir -p .kube || fail_and_exit "can not create kubeconfig directory"

export HOME=/nonexistent

session_url=${KUBEPI_WEBKUBECTL_SESSION_URL:-http://127.0.0.1/kubepi/api/v1/webkubectl/session}
code=$(curl --connect-timeout 5 --max-time 20 -w "%{http_code}" -s -o ~/.kube/config "${session_url}?token=${arg1}")
curl_status=$?

if [[ $curl_status -ne 0 || "${code}" != "200" ]];then
    rm -f .kube/config
    fail_and_exit "download kubeconfig failed (curl: ${curl_status}, http: ${code:-000})"
fi

current_context=$(kubectl config current-context 2>/tmp/kubepi-kubectl.err)
if [[ $? -ne 0 || -z "${current_context}" ]]; then
    cat /tmp/kubepi-kubectl.err
    fail_and_exit "load kubeconfig failed"
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
