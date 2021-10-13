// +build darwin

package server

import (
	"io"
	"os"
	"os/exec"
)

func (e *KubePiSerer) startTty() {
	cmd := "/Users/shenchenyang/go/bin/gotty"
	params := []string{"--permit-write", "bash", "init-kube.sh"}
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
