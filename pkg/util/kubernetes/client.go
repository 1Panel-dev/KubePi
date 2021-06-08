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
	kubeConf := &rest.Config{
		Host:        c.Host,
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	return kubernetes.NewForConfig(kubeConf)

}
