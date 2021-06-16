package kubernetes

import (
	"context"
	"errors"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
)

type Interface interface {
	Ping() error
	Version() (*version.Info, error)
	Client() (*kubernetes.Clientset, error)
}

type Kubernetes struct {
	*v1Cluster.Cluster
}

func NewKubernetes(cluster v1Cluster.Cluster) Interface {
	return &Kubernetes{Cluster: &cluster}
}

func (k *Kubernetes) Client() (*kubernetes.Clientset, error) {
	if k.Spec.Connect.Direction == "forward" {
		kubeConf := &rest.Config{
			Host: k.Spec.Connect.Forward.ApiServer,
		}
		if len(k.CaCertificate.CertData) > 0 {
			kubeConf.CAData = k.CaCertificate.CertData
		} else {
			kubeConf.Insecure = true
		}
		switch k.Spec.Authentication.Mode {
		case strings.ToLower("bearer"):
			kubeConf.BearerToken = k.Spec.Authentication.BearerToken
		case strings.ToLower("certificate"):
			kubeConf.TLSClientConfig.CertData = k.Spec.Authentication.Certificate.CertData
			kubeConf.TLSClientConfig.KeyData=k.Spec.Authentication.Certificate.KeyData
		}
		return kubernetes.NewForConfig(kubeConf)

	}
	return nil, errors.New("")
}

func (k *Kubernetes) Version() (*version.Info, error) {
	client, err := k.Client()
	if err != nil {
		return nil, err
	}
	return client.ServerVersion()

}

func (k *Kubernetes) Ping() error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	_, err = client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	return nil
}
