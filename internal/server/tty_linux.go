// +build linux

package server

func (e *KubePiSerer) startTty() {
	cmd := "gotty"
	params := []string{"--permit-write", "unshare", "--fork", "--pid", "--mount-proc", "--mount", "bash", "init-kube.sh"}
	go func() {
		c := exec.Command(cmd, params...)
		c.Env = append(c.Env, os.Environ()...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			e.logger.Error(err)
		}
	}()
}
