package file

import (
	"github.com/KubeOperator/kubepi/internal/model/v1/file"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/util/kubernetes"
	"github.com/KubeOperator/kubepi/pkg/util/podfile"
)

type Service interface {
	ListFiles(request file.Request) (string, error)
}

type service struct {
	clusterService cluster.Service
}

func NewService() Service {
	return &service{
		clusterService: cluster.NewService(),
	}
}

func (f service) ListFiles(request file.Request) (string, error) {
	var res []string
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
	podcp := podfile.NewPodConfig(request.Namespace, request.PodName, request.ContainerName, k8sConfig, k8sClient)
	return podcp.ListFiles()
}
