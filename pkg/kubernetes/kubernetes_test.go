package kubernetes

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"testing"
)

func TestNewKubernetes(t *testing.T) {
	//cs, err := clientcmd.BuildConfigFromFlags("", "/Users/shenchenyang/.kube/config")
	//if err != nil {
	//	t.Error(err)
	//}

	cs := rest.Config{
		Host: "https://127.0.0.1:51584",
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	cs.BearerToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6Im5sTlV3WWZIcFAwY0lTRF9yVE43YXZyNko1MUhFakQxVktpRGZHRDJrVkUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJkZWZhdWx0LXRva2VuLWZ3NzRnIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImRlZmF1bHQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIzNWExZmFlMS04ZjRiLTQ2ZWEtODA4Mi0yMTg2MjkxODJkOGMiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06ZGVmYXVsdCJ9.MxaaniuIZp5ajqMSD9cplsYidzxr-R8gKceYn4MWmqiPwVXE5XnCAWtXpEuhFCqtaZleZrf6ieQ13LB7FWUC6LWq6oNsbK0UA87w9rBHZL2OKQneEzPj1jS9nGGs03AYAYfBLE5p4zgLRMMW5223Cl_p53sKJjmOgR4F5EntOs9_tbc-qmHh1UE-lf0QB7coXrjp8jrEk_N6hHxM3jfxSwj6MDz54Kogt1w83-P7hwQslYecX1e9KreM0IXKju0bHYrHN2u8iUiZNvGv4HEZaIeSoYWm5uaWePGyhOeeJ-WqP_4bEsp6FKyJBM4yaqnwiwbSSDNZw_xu7uqitTesrQ"

	cls, err := kubernetes.NewForConfig(&cs)
	if err != nil {
		t.Error(err)
	}
	v, _ := cls.ServerVersion()
	fmt.Println(v.String())

}
