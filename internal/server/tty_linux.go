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
		c.Env = append(c.Env, os.Environ()...)
		c.Stdout = io.Discard
		c.Stderr = io.Discard
		if err := c.Run(); err != nil {
			e.logger.Error(err)
		}
	}()
}
