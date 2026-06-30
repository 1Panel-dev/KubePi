package clusteraccess

import (
	"encoding/pem"

	v1Cluster "github.com/1Panel-dev/KubePi/internal/model/v1/cluster"
	"github.com/1Panel-dev/KubePi/internal/service/v1/cluster"
	"github.com/1Panel-dev/KubePi/internal/service/v1/clusterbinding"
	"github.com/1Panel-dev/KubePi/internal/service/v1/common"
	"github.com/1Panel-dev/KubePi/pkg/kubernetes"
	kubernetesClient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type User struct {
	Name            string
	IsAdministrator bool
}

func ConfigForUser(clusterName string, user User) (*rest.Config, *v1Cluster.Cluster, error) {
	c, err := cluster.NewService().Get(clusterName, common.DBOptions{})
	if err != nil {
		return nil, nil, err
	}

	cfg, err := kubernetes.NewKubernetes(c).Config()
	if err != nil {
		return nil, nil, err
	}
	if user.IsAdministrator {
		return cfg, c, nil
	}

	binding, err := clusterbinding.NewService().GetBindingByClusterNameAndUserName(c.Name, user.Name, common.DBOptions{})
	if err != nil {
		return nil, nil, err
	}
	ApplyUserAccessConfig(cfg, c, binding)
	return cfg, c, nil
}
func ApplyUserAccessConfig(cfg *rest.Config, c *v1Cluster.Cluster, binding *v1Cluster.Binding) {
	if len(binding.Certificate) == 0 {
		cfg.Impersonate = rest.ImpersonationConfig{UserName: binding.UserRef}
		return
	}

	cfg.Username = ""
	cfg.Password = ""
	cfg.BearerToken = ""
	cfg.BearerTokenFile = ""
	cfg.AuthProvider = nil
	cfg.ExecProvider = nil
	cfg.Impersonate = rest.ImpersonationConfig{}
	cfg.CertData = binding.Certificate
	cfg.KeyData = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: c.PrivateKey})
}

func ClientForUser(clusterName string, user User) (*kubernetesClient.Clientset, *rest.Config, *v1Cluster.Cluster, error) {
	cfg, c, err := ConfigForUser(clusterName, user)
	if err != nil {
		return nil, nil, nil, err
	}
	client, err := kubernetesClient.NewForConfig(cfg)
	if err != nil {
		return nil, nil, nil, err
	}
	return client, cfg, c, nil
}
