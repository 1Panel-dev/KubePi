package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	Host  string
	Token string
}

func NewKubernetesClient(c *Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(NewClusterConfig(c))
}

func NewClusterConfig(c *Config) *rest.Config {
	return &rest.Config{
		Host:        c.Host,
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
}
