// +build darwin

package server

import (
	"os"
	"os/exec"
)

func (e *KubePiSerer) startTty() {
	cmd := "/Users/shenchenyang/go/bin/gotty"
	params := []string{"--permit-write", "bash", "init-kube.sh"}
	go func() {
		c := exec.Command(cmd, params...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			e.logger.Error(err)
		}
	}()

}
