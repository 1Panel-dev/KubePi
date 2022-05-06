package podfile

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type PodCp struct {
	K8sClient     *kubernetes.Clientset
	RESTConfig    *rest.Config
	Namespace     string
	PodName       string
	ContainerName string
	Command       []string
	Stdin         io.Reader
	Stdout        io.Writer
	Stderr        io.Writer
	Tty           bool
	NoPreserve    bool
}

func NewPodConfig(namespace, podName, containerName string, restConfig *rest.Config, k8sClient *kubernetes.Clientset) *PodCp {
	return &PodCp{
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		RESTConfig:    restConfig,
		K8sClient:     k8sClient,
	}
}

type ActionType string

const Upload ActionType = "Upload"
const Download ActionType = "Download"

func (p *PodCp) CopyToPod(srcPath, destPath string) error {
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		cmdutil.CheckErr(makeTar(srcPath, destPath, writer))
	}()
	p.Tty = false
	p.NoPreserve = false
	p.Stdin = reader
	p.Stdout = os.Stdout
	if p.NoPreserve {
		p.Command = []string{"tar", "--no-same-permissions", "--no-same-owner", "-xmf", "-"}
	} else {
		p.Command = []string{"tar", "-xmf", "-"}
	}
	var stderr bytes.Buffer
	p.Stderr = &stderr
	err := p.Exec(Upload)
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

func (p *PodCp) CopyFromPod(filePath string, destPath string) error {
	reader, outStream := io.Pipe()

	p.Command = []string{"tar", "cf", "-", filePath}
	p.Stdin = os.Stdin
	p.Stdout = outStream
	p.Stderr = os.Stderr

	err := p.Exec(Download)
	if err != nil {
		return err
	}
	prefix := getPrefix(filePath)
	prefix = path.Clean(prefix)
	prefix = stripPathShortcuts(prefix)
	err = unTarAll(reader, destPath, prefix)
	return err
}

func (p *PodCp) ListFiles() (string, error) {
	p.Command = []string{"ls", "-lQ", "--color=never", "--full-time", "/"}
	var stdout, stderr bytes.Buffer
	p.Stdout = &stdout
	p.Stderr = &stderr
	p.Tty = false
	err := p.Exec(Upload)
	if err != nil {
		return "", err
	}
	if stderr.String() != "" {
		err = errors.New(stderr.String())
	}

	return stdout.String(), err
}

func (p *PodCp) Exec(actionType ActionType) error {
	req := p.K8sClient.CoreV1().RESTClient().Get().
		Resource("pods").
		Name(p.PodName).
		Namespace(p.Namespace).
		SubResource("exec").
		VersionedParams(&coreV1.PodExecOptions{
			Command:   p.Command,
			Container: p.ContainerName,
			Stdin:     p.Stdin != nil,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(p.RESTConfig, "POST", req.URL())
	if err != nil {
		return err
	}

	if actionType == Upload {
		err = p.stream(exec)
	}
	if actionType == Download {
		go func() {
			p.stream(exec)
		}()
	}
	return err
}

func (p *PodCp) stream(exec remotecommand.Executor) error {
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  p.Stdin,
		Stdout: p.Stdout,
		Stderr: p.Stderr,
		Tty:    p.Tty,
	})
}

func makeTar(srcPath, destPath string, writer io.Writer) error {
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

func unTarAll(reader io.Reader, destDir, prefix string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}

		mode := header.FileInfo().Mode()
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])

		baseName := filepath.Dir(destFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFileName, 0755); err != nil {
				return err
			}
			continue
		}

		evaledPath, err := filepath.EvalSymlinks(baseName)
		if err != nil {
			return err
		}

		if mode&os.ModeSymlink != 0 {
			linkname := header.Linkname

			if !filepath.IsAbs(linkname) {
				_ = filepath.Join(evaledPath, linkname)
			}

			if err := os.Symlink(linkname, destFileName); err != nil {
				return err
			}
		} else {
			outFile, err := os.Create(destFileName)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func getPrefix(file string) string {
	return strings.TrimLeft(file, "/")
}

func stripPathShortcuts(p string) string {

	newPath := path.Clean(p)
	trimmed := strings.TrimPrefix(newPath, "../")

	for trimmed != newPath {
		newPath = trimmed
		trimmed = strings.TrimPrefix(newPath, "../")
	}

	if newPath == "." || newPath == ".." {
		newPath = ""
	}

	if len(newPath) > 0 && string(newPath[0]) == "/" {
		return newPath[1:]
	}

	return newPath
}
