package main

import (
	"context"
	"fmt"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"strings"
)

func main() {
	opts := sockjs.DefaultOptions
	opts.Websocket = true
	handler := sockjs.NewHandler("/echo", opts, echoHandler)
	http.Handle("/echo/", handler)
	log.Println("Server started on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func echoHandler(session sockjs.Session) {
	cs := rest.Config{
		Host: "https://127.0.0.1:55155",
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	cs.BearerToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6Im5sTlV3WWZIcFAwY0lTRF9yVE43YXZyNko1MUhFakQxVktpRGZHRDJrVkUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJkZWZhdWx0LXRva2VuLWZ3NzRnIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImRlZmF1bHQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIzNWExZmFlMS04ZjRiLTQ2ZWEtODA4Mi0yMTg2MjkxODJkOGMiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06ZGVmYXVsdCJ9.MxaaniuIZp5ajqMSD9cplsYidzxr-R8gKceYn4MWmqiPwVXE5XnCAWtXpEuhFCqtaZleZrf6ieQ13LB7FWUC6LWq6oNsbK0UA87w9rBHZL2OKQneEzPj1jS9nGGs03AYAYfBLE5p4zgLRMMW5223Cl_p53sKJjmOgR4F5EntOs9_tbc-qmHh1UE-lf0QB7coXrjp8jrEk_N6hHxM3jfxSwj6MDz54Kogt1w83-P7hwQslYecX1e9KreM0IXKju0bHYrHN2u8iUiZNvGv4HEZaIeSoYWm5uaWePGyhOeeJ-WqP_4bEsp6FKyJBM4yaqnwiwbSSDNZw_xu7uqitTesrQ"
	cls, err := kubernetes.NewForConfig(&cs)
	if err != nil {
		return
	}
	var lines int64
	lines = 20
	reader, err := cls.CoreV1().Pods("kube-system").GetLogs("kube-apiserver-minikube", &v1.PodLogOptions{
		Follow:    false,
		TailLines: &lines,
	}).Stream(context.TODO())
	if err != nil {
		return
	}
	log.Println("new sockjs session established")
	for {
		buf := make([]byte, 2048)
		numBytes, err := reader.Read(buf)
		if numBytes > 0 {
			message := string(buf[:numBytes])
			for _, val := range strings.Split(message, "\n") {
				if strings.TrimSpace(val) == "" {
					continue
				}
				_ = session.Send(strings.TrimSpace(val) + "\r\n")
			}
			fmt.Println(message)
		}
		if err != nil {
			break
		}
	}
	log.Println("sockjs session closed")
}
