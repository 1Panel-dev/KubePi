//go:build linux
// +build linux

package server

import (
	"io"
	"os"
	"os/exec"
)

func (e *KubePiServer) startTty() {
	cmd := "gotty"
	params := []string{"--permit-write", "unshare", "--fork", "--pid", "--mount-proc", "--mount", "bash", "init-kube.sh"}
	go func() {
		c := exec.Command(cmd, params...)
		c.Env = append(os.Environ(), "KUBEPI_WEBKUBECTL_SESSION_URL="+e.localWebkubectlSessionURL())
		c.Stdout = io.Discard
		c.Stderr = io.Discard
		if err := c.Run(); err != nil {
			e.logger.Error(err)
		}
	}()
}
