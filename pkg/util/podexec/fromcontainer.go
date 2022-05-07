package podexec

import (
	"bytes"
	"fmt"
)

// FromPodContainer 从pod内拷贝到io.Writer
func (p *PodExec) FromPodContainer(dest []string, style string) error {
	switch style {
	case "tar":
		p.Command = append([]string{"tar", "cf", "-"}, dest...)
	case "zip":
		p.Command = append([]string{"/kotools", "zip"}, dest...)
	default:
		p.Command = append([]string{"tar", "cf", "-"}, dest...)
	}
	p.Tty = false
	var stderr bytes.Buffer
	p.Stderr = &stderr
	err := p.Exec()
	if err != nil {
		if len(stderr.Bytes()) != 0 {
			return fmt.Errorf("STDERR: " + stderr.String())
		}
		return fmt.Errorf(err.Error(), stderr)
	}
	return nil
}
