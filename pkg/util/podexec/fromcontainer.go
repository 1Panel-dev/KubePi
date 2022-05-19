package podexec

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
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
	err := p.Exec(Exec)
	if err != nil {
		if len(stderr.Bytes()) != 0 {
			return fmt.Errorf("STDERR: " + stderr.String())
		}
		return fmt.Errorf(err.Error(), stderr)
	}
	return nil
}

func (p *PodExec) CopyFromPod(filePath string, destPath string) error {
	reader, outStream := io.Pipe()

	p.Command = []string{"tar", "cf", "-", filePath}
	p.Stdin = os.Stdin
	p.Stdout = outStream
	p.Stderr = os.Stderr
	p.NoPreserve = true

	err := p.Exec(Download)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	r := bufio.NewReader(reader)
	w := bufio.NewWriter(file)
	size := 4 * 1024
	buf := make([]byte, 4*1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		_, err = w.Write(buf[:n])
		if err != nil {
			return err
		}
		if n < size {
			break
		}
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return err
}
