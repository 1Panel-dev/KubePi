package podbase

import (
	"bytes"
	"context"
	"fmt"
	"github.com/KubeOperator/kubepi/pkg/util/podexec"
	"github.com/sirupsen/logrus"
	"io"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

//const KotoolsPath = "/kubepi/kotools"

const KotoolsPath = "/Users/zk.wang/go/src/github.com/KubeOperator/kotools/utils/binary"

type PodBase struct {
	Namespace  string
	PodName    string
	Container  string
	K8sClient  *kubernetes.Clientset
	RestClient *rest.Config
}

func (p *PodBase) NewPodExec() *podexec.PodExec {
	return podexec.NewPodExec(p.Namespace, p.PodName, p.Container, p.RestClient, p.K8sClient)
}

func (p *PodBase) PodInfo() (*coreV1.Pod, error) {
	pod, err := p.K8sClient.CoreV1().Pods(p.Namespace).
		Get(context.TODO(), p.PodName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if pod.Status.Phase == coreV1.PodSucceeded || pod.Status.Phase == coreV1.PodFailed {
		return nil, fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
	}
	return pod, nil
}

func (p *PodBase) OsAndArch(nodeName string) (osType string, arch string) {
	// get pod system arch and type
	node, err := p.K8sClient.CoreV1().Nodes().
		Get(context.TODO(), nodeName, metaV1.GetOptions{})
	if err == nil {
		var ok bool
		osType, ok = node.Labels["beta.kubernetes.io/os"]
		if !ok {
			osType, ok = node.Labels["kubernetes.io/os"]
			if !ok {
				osType = "linux"
			}
		}
		arch, ok = node.Labels["beta.kubernetes.io/arch"]
		if !ok {
			arch, ok = node.Labels["kubernetes.io/arch"]
			if !ok {
				arch = "amd64"
			}
		}
	}
	return
}

func (p *PodBase) Exec(stdin io.Reader, command ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	exec := p.NewPodExec()
	exec.Command = command
	exec.Tty = false
	exec.Stdin = stdin
	exec.Stdout = &stdout
	exec.Stderr = &stderr
	err := exec.Exec()
	if err != nil {
		if len(stderr.Bytes()) != 0 {
			return nil, fmt.Errorf("%s %s", err.Error(), stderr.String())
		}
		return nil, err
	}
	if len(stderr.Bytes()) != 0 {
		return nil, fmt.Errorf(stderr.String())
	}
	return stdout.Bytes(), nil
}

func (p *PodBase) InstallKFTools() error {
	pod, err := p.PodInfo()
	if err != nil {
		return err
	}
	osType, arch := p.OsAndArch(pod.Spec.NodeName)
	kfToolsPath := fmt.Sprintf(KotoolsPath+"/kotools_%s_%s", osType, arch)
	if osType == "windows" {
		kfToolsPath += ".exe"
	}
	exec := p.NewPodExec()
	err = exec.CopyToPod(kfToolsPath, "/kotools")
	if err != nil {
		return err
	}
	if osType != "windows" {
		chmodCmd := []string{"chmod", "+x", "/kotools"}
		exec.Command = chmodCmd
		var stderr bytes.Buffer
		exec.Stderr = &stderr
		err = exec.Exec()
		if err != nil {
			return fmt.Errorf(err.Error(), stderr)
		}
	}
	return nil
}

func TarKFTools(name string, writer io.Writer) error {
	//tw := tar.NewWriter(writer)
	// 如果关闭失败会造成tar包不完整
	//defer tw.Close()
	f, err := os.Open(name)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer f.Close()
	//fi, err := f.Stat()
	//if err != nil {
	//	logrus.Error(err)
	//	return err
	//}
	//hdr, err := tar.FileInfoHeader(fi, name)
	//if err != nil {
	//	logrus.Error(err)
	//	return err
	//}
	//hdr.Name = "kotools"
	//// 将tar的文件信息hdr写入到tw
	//err = tw.WriteHeader(hdr)
	//if err != nil {
	//	logrus.Error(err)
	//	return err
	//}
	// 将文件数据写入
	_, err = io.Copy(writer, f)

	return err
}
