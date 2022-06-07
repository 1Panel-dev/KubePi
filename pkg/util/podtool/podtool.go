package podtool

import (
	"bytes"
	"io"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type PodTool struct {
	Namespace     string
	PodName       string
	ContainerName string
	K8sClient     *kubernetes.Clientset
	RestClient    *rest.Config
	ExecConfig    ExecConfig
}

type ExecConfig struct {
	Command    []string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	Tty        bool
	NoPreserve bool
}

func (p *PodTool) ExecCommand(commands []string) ([]byte, error) {
	var stdout bytes.Buffer
	p.ExecConfig.Stdout = &stdout
	p.ExecConfig.Command = commands
	err := p.Exec(Exec)
	if err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}

type ActionType string

const Exec ActionType = "Exec"
const Download ActionType = "Download"

// Exec 在给定容器中执行命令
func (p *PodTool) Exec(actionType ActionType) error {
	config := p.ExecConfig
	req := p.K8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(p.PodName).
		Namespace(p.Namespace).
		SubResource("exec").
		VersionedParams(&coreV1.PodExecOptions{
			Command:   config.Command,
			Container: p.ContainerName,
			Stdin:     config.Stdin != nil,
			Stdout:    config.Stdout != nil,
			Stderr:    config.Stderr != nil,
			TTY:       config.Tty,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(p.RestClient, "POST", req.URL())
	if err != nil {
		return err
	}
	if actionType == Download {
		go func() {
			_ = p.stream(exec)
		}()
	} else {
		err = p.stream(exec)
	}
	return err
}

func (p *PodTool) stream(exec remotecommand.Executor) error {
	config := p.ExecConfig
	var sizeQueue remotecommand.TerminalSizeQueue
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:             config.Stdin,
		Stdout:            config.Stdout,
		Stderr:            config.Stderr,
		Tty:               config.Tty,
		TerminalSizeQueue: sizeQueue,
	})
}
