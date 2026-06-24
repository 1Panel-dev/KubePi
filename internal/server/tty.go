package server

import (
	"net"
	"strconv"
	"strings"
)

func (e *KubePiServer) localWebkubectlSessionURL() string {
	host := strings.TrimSpace(e.config.Spec.Server.Bind.Host)
	switch host {
	case "", "0.0.0.0", "::", "[::]":
		host = "127.0.0.1"
	}
	return "http://" + net.JoinHostPort(host, strconv.Itoa(e.config.Spec.Server.Bind.Port)) + "/kubepi/api/v1/webkubectl/session"
}
