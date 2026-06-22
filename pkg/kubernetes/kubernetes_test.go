package kubernetes

import (
	"os"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func TestNewKubernetes(t *testing.T) {
	host := os.Getenv("KUBERNETES_TEST_HOST")
	token := os.Getenv("KUBERNETES_TEST_TOKEN")
	if host == "" || token == "" {
		t.Skip("set KUBERNETES_TEST_HOST and KUBERNETES_TEST_TOKEN to run live Kubernetes client test")
	}

	cs := rest.Config{
		Host: host,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	cs.BearerToken = token

	cls, err := kubernetes.NewForConfig(&cs)
	if err != nil {
		t.Fatal(err)
	}
	v, err := cls.ServerVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v.String())
}
