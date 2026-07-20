#!/bin/bash
set -u

start_shell() {
    export BASH_SILENCE_DEPRECATION_WARNING=1
    exec /bin/bash
}

fail_and_exit() {
    echo "$1" >&2
    exit 1
}

arg1=${1:-}
if [[ -z "${arg1}" ]]; then
    fail_and_exit "missing terminal session token"
fi

terminal_home=$(mktemp -d "${TMPDIR:-/tmp}/kubepi-terminal.XXXXXX") || fail_and_exit "can not create terminal home"
export HOME="${terminal_home}"
mkdir -p ~/.kube || fail_and_exit "can not create kubeconfig directory"

session_url=${KUBEPI_WEBKUBECTL_SESSION_URL:-http://127.0.0.1/kubepi/api/v1/webkubectl/session}
code=$(curl --connect-timeout 5 --max-time 20 -w "%{http_code}" -s -o ~/.kube/config "${session_url}?token=${arg1}")
curl_status=$?

if [[ $curl_status -ne 0 || "${code}" != "200" ]];then
    rm -f ~/.kube/config
    fail_and_exit "download kubeconfig failed (curl: ${curl_status}, http: ${code:-000})"
fi

current_context=$(kubectl config current-context 2>/tmp/kubepi-kubectl.err)
if [[ $? -ne 0 || -z "${current_context}" ]]; then
    cat /tmp/kubepi-kubectl.err
    fail_and_exit "load kubeconfig failed"
fi
cluster=${current_context%@*}
username=${current_context#*@}

echo "Welcome to kubepi"
echo "Current cluster is ${cluster}"
echo "Current user is ${username}"

start_shell
