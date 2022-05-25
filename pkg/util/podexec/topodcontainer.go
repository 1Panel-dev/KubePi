package podexec

import (
	"bytes"
	"fmt"
	"io"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"path"
	"strings"
)

func (p *PodExec) CopyToPod(srcPath, destPath string) error {
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		cmdutil.CheckErr(MakeTar(srcPath, destPath, writer))
	}()
	p.Tty = false
	p.NoPreserve = false
	p.Stdin = reader
	if p.NoPreserve {
		p.Command = []string{"tar", "--no-same-permissions", "--no-same-owner", "-xmf", "-"}
	} else {
		p.Command = []string{"tar", "-xmf", "-"}
	}
	var stdout bytes.Buffer
	p.Stdout = &stdout
	dest := path.Dir(destPath)
	p.Command = append(p.Command, "-C", dest)
	var stderr bytes.Buffer
	p.Stderr = &stderr
	err := p.Exec(Exec)
	if err != nil {
		return fmt.Errorf(err.Error(), stdout.String())
	}
	if len(stderr.Bytes()) != 0 {
		for _, line := range strings.Split(stderr.String(), "\n") {
			if len(strings.TrimSpace(line)) == 0 {
				continue
			}
			if !strings.Contains(strings.ToLower(line), "removing") {
				return fmt.Errorf(line)
			}
		}
	}
	return nil
}
