package proxy

import (
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/clusterbinding"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Handler struct {
	clusterService        cluster.Service
	clusterBindingService clusterbinding.Service
}

func NewHandler() *Handler {
	return &Handler{
		clusterService:        cluster.NewService(),
		clusterBindingService: clusterbinding.NewService(),
	}
}

type NamespaceResourceContainer struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []interface{} `json:"items"`
	Namespaces      []string      `json:"namespaces"`
}

func (h *Handler) KubernetesAPIProxy() iris.Handler {
	return func(ctx *context.Context) {
		// 解析参数
		name := ctx.Params().GetString("name")
		proxyPath := ensureProxyPathValid(ctx.Params().GetString("p"))
		namespace := ctx.URLParam("namespace")
		requestMethod := ctx.Request().Method
		// 获取当亲集群
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		// 获取session
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		// 生成transport
		ts, err := h.generateTLSTransport(c, profile.Name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		// 生成httpClient
		httpClient := http.Client{Transport: ts}
		k := kubernetes.NewKubernetes(*c)
		// 判断资源类型是否是namespace级别的
		resourceName, err := parseResourceName(proxyPath)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		namespaced, err := k.IsNamespacedResource(resourceName)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		canVisitAll, err := k.CanVisitAllNamespace(profile.Name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		//
		apiUrl, err := url.Parse(fmt.Sprintf("%s%s", c.Spec.Connect.Forward.ApiServer, proxyPath))
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		apiUrl.RawQuery = ctx.Request().URL.RawQuery
		if http.MethodGet == requestMethod && namespace == "" && namespaced && !canVisitAll {
			// 调用多namespace 逻辑
			allowedNamespaces, err := k.GetUserNamespaceNames(profile.Name)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err)
				return
			}
			resp, err := fetchMultiNamespaceResource(&httpClient, allowedNamespaces, *apiUrl)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err)
				return
			}
			_, _ = ctx.JSON(resp)
			return
		}
		if http.MethodGet == requestMethod && namespaced && namespace != "" {
			apiUrl.Path = addUrlNamespace(apiUrl.Path, namespace)
		}
		req, err := http.NewRequest(ctx.Request().Method, apiUrl.String(), ctx.Request().Body)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusForbidden {
				resp.StatusCode = http.StatusInternalServerError
			}
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		rawResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		_, _ = ctx.Write(rawResp)
	}
}

func fetchMultiNamespaceResource(client *http.Client, namespaces []string, apiUrl url.URL) (*NamespaceResourceContainer, error) {
	wg := &sync.WaitGroup{}
	var mergedContainer NamespaceResourceContainer
	var responses []*http.Response
	var es []error
	for i := range namespaces {
		wg.Add(1)
		ns := namespaces[i]
		go func() {
			newUrl := apiUrl
			newUrl.Path = addUrlNamespace(apiUrl.Path, ns)
			resp, err := client.Get(newUrl.String())
			if err != nil {
				es = append(es, err)
				wg.Done()
				return
			}
			responses = append(responses, resp)
			wg.Done()
		}()

	}
	wg.Wait()
	var forbidden int
	var forbiddenMessage []string
	for i := range responses {
		r := responses[i]
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		if r.StatusCode != http.StatusOK {
			if r.StatusCode == http.StatusForbidden {
				forbidden++
				forbiddenMessage = append(forbiddenMessage, string(body))
				continue
			} else {
				return nil, errors.New(string(body))
			}
		}
		var nc NamespaceResourceContainer
		if err := json.Unmarshal(body, &nc); err != nil {
			return nil, err
		}
		mergedContainer.TypeMeta = nc.TypeMeta
		mergedContainer.ListMeta = nc.ListMeta
		mergedContainer.Items = append(mergedContainer.Items, nc.Items...)
	}
	if len(namespaces) == 1 && forbidden == 1 {
		return nil, errors.New(strings.Join(forbiddenMessage, ""))
	}
	return &mergedContainer, nil
}

func (h *Handler) generateTLSTransport(c *v1Cluster.Cluster, username string) (http.RoundTripper, error) {
	binding, err := h.clusterBindingService.GetBindingByClusterNameAndUserName(c.Name, username, common.DBOptions{})
	if err != nil {
		return nil, err
	}
	kubeConf := &rest.Config{
		Host: c.Spec.Connect.Forward.ApiServer,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
			CertData: binding.Certificate,
			KeyData:  pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: c.PrivateKey}),
		},
	}
	ts, err := rest.TransportFor(kubeConf)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func ensureProxyPathValid(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = "/" + path
	}
	return path
}

func parseResourceName(path string) (string, error) {
	ss := strings.Split(path, "/")
	if len(ss) > 0 {
		return ss[len(ss)-1], nil
	}
	return "", fmt.Errorf("cant not get resource name from url %s", path)
}

func addUrlNamespace(path string, ns string) string {
	ss := strings.Split(path, "/")
	resourceName := ss[len(ss)-1]
	per := ss[:len(ss)-1]
	namespacedSs := append(per, "namespaces", ns, resourceName)
	return strings.Join(namespacedSs, "/")
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/proxy")
	sp.Any("/:name/k8s/{p:path}", handler.KubernetesAPIProxy())
}
