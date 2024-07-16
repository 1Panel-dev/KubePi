package podtool

import (
	"bufio"
	"io"
	"os"
)
func (p *PodTool) CopyFileFromPod(filePath string, destPath string) error {
	reader, outStream := io.Pipe()

	p.ExecConfig = ExecConfig{
		Command:    []string{"cat", filePath},
		Stdin:      os.Stdin,
		Stdout:     outStream,
		Stderr:     os.Stderr,
		NoPreserve: true,
	}

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
func (p *PodTool) CopyFolderFromPod(filePath string, destPath string) error {
	reader, outStream := io.Pipe()

	p.ExecConfig = ExecConfig{
		Command:    []string{"tar", "cf", "-", filePath},
		Stdin:      os.Stdin,
		Stdout:     outStream,
		Stderr:     os.Stderr,
		NoPreserve: true,
	}

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
