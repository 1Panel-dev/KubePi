// +build linux

package server

func (e *KubePiSerer) startTty() {
	cmd := "gotty"
	params := []string{"--permit-write", "bash", "unshare", "--fork", "--pid", "--mount-proc", "--mount", "init-kube.sh"}
	go func() {
		c := exec.Command(cmd, params...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			e.logger.Error(err)
		}
	}()
}
