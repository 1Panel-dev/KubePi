package file

import (
	"fmt"
	"github.com/KubeOperator/kubepi/internal/model/v1/file"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	kubeClient "github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/KubeOperator/kubepi/pkg/util/podtool"
	"github.com/sirupsen/logrus"
	"io"
	"k8s.io/client-go/kubernetes"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Service interface {
	ListFiles(request file.Request) ([]podtool.File, error)
	DownloadFile(request file.Request) (string, error)
	DownloadFolder(request file.Request) (string, error)
	UploadFile(request file.Request) error
	ExecNewCommand(request file.Request) ([]byte, error)
	EditFile(request file.Request) error
	CatFile(request file.Request) ([]byte, error)
}

type service struct {
	clusterService cluster.Service
}

func NewService() Service {
	return &service{
		clusterService: cluster.NewService(),
	}
}

func (f service) ExecNewCommand(request file.Request) ([]byte, error) {
	pt, err := f.GetPodTool(request)
	if err != nil {
		return nil, err
	}
	return pt.ExecCommand(request.Commands)
}

func (f service) GetPodTool(request file.Request) (podtool.PodTool, error) {
	var pt podtool.PodTool
	clu, err := f.clusterService.Get(request.Cluster, common.DBOptions{})
	if err != nil {
		return pt, err
	}
	client := kubeClient.Kubernetes{clu}
	config, err := client.Config()
	if err != nil {
		return pt, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return pt, err
	}
	pt = podtool.PodTool{
		Namespace:     request.Namespace,
		PodName:       request.PodName,
		ContainerName: request.ContainerName,
		K8sClient:     clientset,
		RestClient:    config,
		ExecConfig: podtool.ExecConfig{
			Stdin: request.Stdin,
		},
	}
	return pt, nil
}

func (f service) ListFiles(request file.Request) ([]podtool.File, error) {
	pt, err := f.GetPodTool(request)
	if err != nil {
		return nil, err
	}
	return pt.ListFiles(request.Path)
}

func (f service) EditFile(request file.Request) error {
	pt, err := f.GetPodTool(request)
	if err != nil {
		return err
	}
	return pt.EditFile(request.Path, request.Content)
}

func (f service) DownloadFile(request file.Request) (string, error) {

	var fileP string
	pt, err := f.GetPodTool(request)
	if err != nil {
		return fileP, err
	}
	fileNameWithSuffix := path.Base(request.Path)
	fileP = filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
	err = os.MkdirAll(fileP, os.ModePerm)
	if err != nil {
		return "", err
	}
	fileP = filepath.Join(fileP, fileNameWithSuffix)
	err = pt.CopyFileFromPod(request.Path, fileP)
	if err != nil {
		return "", err
	}

	return fileP, nil
}
func (f service) DownloadFolder(request file.Request) (string, error) {

	var fileP string
	pt, err := f.GetPodTool(request)
	if err != nil {
		return fileP, err
	}
	fileNameWithSuffix := path.Base(request.Path)
	fileType := path.Ext(fileNameWithSuffix)
	fileName := strings.TrimSuffix(fileNameWithSuffix, fileType)
	fileP = filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
	err = os.MkdirAll(fileP, os.ModePerm)
	if err != nil {
		return "", err
	}
	fileP = filepath.Join(fileP, fileName+".tar")
	err = pt.CopyFolderFromPod(request.Path, fileP)
	if err != nil {
		return "", err
	}

	return fileP, nil
}

func (f service) UploadFile(request file.Request) error {
	pt, err := f.GetPodTool(request)
	if err != nil {
		return err
	}
	reader, writer := io.Pipe()
	pt.ExecConfig.Stdin = reader
	go func() {
		defer func() {
			_ = writer.Close()
		}()
		tarFile, err := os.Open(request.FilePath)
		if err != nil {
			logrus.Error(err)
			return
		}
		_, err = io.Copy(writer, tarFile)
		if err != nil {
			logrus.Error(err)
		}
	}()
	return pt.CopyToContainer(request.Path)
}

func (f service) CatFile(request file.Request) ([]byte, error) {
	pt, err := f.GetPodTool(request)
	if err != nil {
		return nil, err
	}
	return pt.CatFile(request.Path)
}
