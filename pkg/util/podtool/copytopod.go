package podtool

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"os"
	"path"
	"path/filepath"
)

func (p *PodTool) CopyToContainer(destPath string) error {
	if p.ExecConfig.NoPreserve {
		p.ExecConfig.Command = []string{"tar", "--no-same-permissions", "--no-same-owner", "-xmf", "-"}
	} else {
		p.ExecConfig.Command = []string{"tar", "-xmf", "-"}
	}
	if len(destPath) > 0 {
		p.ExecConfig.Command = append(p.ExecConfig.Command, "-C", destPath)
	}
	p.ExecConfig.Tty = false
	var stderr bytes.Buffer
	p.ExecConfig.Stderr = &stderr
	err := p.Exec(Exec)
	var stdout bytes.Buffer
	p.ExecConfig.Stdout = &stdout
	if err != nil {
		result := ""
		if len(stdout.Bytes()) != 0 {
			result = stdout.String()
		}
		if len(stderr.Bytes()) != 0 {
			result = stderr.String()
		}
		return fmt.Errorf(err.Error(), result)
	}
	return nil
}

func (p *PodTool) CopyToPod(srcPath, destPath string) error {
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		cmdutil.CheckErr(MakeTar(srcPath, destPath, writer))
	}()
	config := ExecConfig{
		Tty:        false,
		NoPreserve: false,
		Stdin:      reader,
	}

	if config.NoPreserve {
		config.Command = []string{"tar", "--no-same-permissions", "--no-same-owner", "-xmf", "-"}
	} else {
		config.Command = []string{"tar", "-xmf", "-"}
	}
	var stdout bytes.Buffer
	config.Stdout = &stdout
	dest := path.Dir(destPath)
	config.Command = append(config.Command, "-C", dest)
	var stderr bytes.Buffer
	config.Stderr = &stderr
	p.ExecConfig = config
	err := p.Exec(Exec)
	if err != nil {
		result := ""
		if len(stdout.Bytes()) != 0 {
			result = stdout.String()
		}
		if len(stderr.Bytes()) != 0 {
			result = stderr.String()
		}
		return fmt.Errorf(err.Error(), result)
	}
	return nil
}

func MakeTar(srcPath, destPath string, writer io.Writer) error {
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	srcPath = path.Clean(srcPath)
	destPath = path.Clean(destPath)
	return recursiveTar(path.Dir(srcPath), path.Base(srcPath), path.Dir(destPath), path.Base(destPath), tarWriter)
}

func recursiveTar(srcBase, srcFile, destBase, destFile string, tw *tar.Writer) error {
	srcPath := path.Join(srcBase, srcFile)
	matchedPaths, err := filepath.Glob(srcPath)
	if err != nil {
		return err
	}
	for _, fpath := range matchedPaths {
		stat, err := os.Lstat(fpath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			files, err := ioutil.ReadDir(fpath)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				//case empty directory
				hdr, _ := tar.FileInfoHeader(stat, fpath)
				hdr.Name = destFile
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
			}
			for _, f := range files {
				if err := recursiveTar(srcBase, path.Join(srcFile, f.Name()), destBase, path.Join(destFile, f.Name()), tw); err != nil {
					return err
				}
			}
			return nil
		} else if stat.Mode()&os.ModeSymlink != 0 {
			//case soft link
			hdr, _ := tar.FileInfoHeader(stat, fpath)
			target, err := os.Readlink(fpath)
			if err != nil {
				return err
			}

			hdr.Linkname = target
			hdr.Name = destFile
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
			//case regular file or other file type like pipe
			hdr, err := tar.FileInfoHeader(stat, fpath)
			if err != nil {
				return err
			}
			hdr.Name = destFile

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
			return f.Close()
		}
	}
	return nil
}
