package podexec

import (
	"bytes"
	"fmt"
	"io"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"os"
	"strings"
)

// ToPodContainer 从io.Writer拷贝到pod
// tar file ---> io.Writer ---> kubernetes api ---> io.Reader ---> unTar file
func (p *PodExec) ToPodContainer(destPath string) error {
	if p.NoPreserve {
		p.Command = []string{"tar", "--no-same-permissions", "--no-same-owner", "-xmf", "-"}
	} else {
		p.Command = []string{"tar", "-xmf", "-"}
	}
	//destPath = path.Dir(destPath)
	if len(destPath) > 0 {
		p.Command = append(p.Command, "-C", destPath)
	}
	p.Tty = false
	var stderr bytes.Buffer
	p.Stderr = &stderr
	err := p.Exec()
	if err != nil {
		return fmt.Errorf(err.Error(), stderr)
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

func (p *PodExec) CopyToPod(srcPath, destPath string) error {
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		cmdutil.CheckErr(MakeTar(srcPath, destPath, writer))
	}()
	p.Tty = false
	p.NoPreserve = true
	p.Stdin = reader
	p.Stdout = os.Stdout
	if p.NoPreserve {
		p.Command = []string{"tar", "--no-same-permissions", "--no-same-owner", "-xmf", "-"}
	} else {
		p.Command = []string{"tar", "-xmf", "-"}
	}
	var stderr bytes.Buffer
	p.Stderr = &stderr
	err := p.Exec()
	if err != nil {
		return fmt.Errorf(err.Error(), p.Stderr)
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
