package file

import (
	"encoding/json"
	"github.com/KubeOperator/kubepi/internal/model/v1/file"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/util"
	"github.com/KubeOperator/kubepi/pkg/util/kubernetes"
	"github.com/KubeOperator/kubepi/pkg/util/podbase"
	"github.com/sirupsen/logrus"
	"strings"
)

type Service interface {
	ListFiles(request file.Request) ([]util.File, error)
}

type service struct {
	clusterService cluster.Service
}

func NewService() Service {
	return &service{
		clusterService: cluster.NewService(),
	}
}

func (f service) ListFiles(request file.Request) ([]util.File, error) {

	request.Commands = []string{"/kotools", "ls", request.Path}
	var res []util.File
	bs, err := f.fileBrowser(request)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bs, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (f service) fileBrowser(request file.Request) (res []byte, err error) {
	clu, err := f.clusterService.Get(request.Cluster, common.DBOptions{})
	if err != nil {
		return nil, err
	}
	config := &kubernetes.Config{
		Host:  clu.Spec.Connect.Forward.ApiServer,
		Token: clu.Spec.Authentication.BearerToken,
	}
	k8sConfig := kubernetes.NewClusterConfig(config)
	k8sClient, err := kubernetes.NewKubernetesClient(config)
	if err != nil {
		return nil, err
	}
	pb := podbase.PodBase{
		Namespace:  request.Namespace,
		PodName:    request.PodName,
		Container:  request.ContainerName,
		K8sClient:  k8sClient,
		RestClient: k8sConfig,
	}
	res, err = pb.Exec(request.Stdin, request.Commands...)
	if err != nil {
		if strings.Contains(err.Error(), "kotools") ||
			err.Error() == "command terminated with exit code 126" {
			err = pb.InstallKFTools()
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			res, err = pb.Exec(request.Stdin, request.Commands...)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
		} else {
			logrus.Error(err)
			return nil, err
		}
	}
	return res, err
}
