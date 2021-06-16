package proxy

import (
	"github.com/kataras/iris/v12/context"
)

func KubernetesClientProxy(ctx *context.Context) {
	//clusterName := ctx.Params().Get("cluster_name")
	//proxyPath := ctx.Params().Get("p")
	//cluster, err := clusterService.Get(clusterName)
	//if err != nil {
	//	_, _ = ctx.JSON(iris.StatusInternalServerError)
	//	return
	//}
	//u, err := url.Parse(cluster.ApiServer)
	//if err != nil {
	//	_, _ = ctx.JSON(iris.StatusInternalServerError)
	//	return
	//}
	//proxy := httputil.NewSingleHostReverseProxy(u)
	//proxy.Transport = &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//token := fmt.Sprintf("%s %s", keyPrefix, cluster.Token)
	//ctx.Request().Header.Add(AuthorizationHeader, token)
	//ctx.Request().URL.Path = proxyPath
	//proxy.ModifyResponse = func(response *http.Response) error {
	//	if response.StatusCode == http.StatusUnauthorized {
	//		response.StatusCode = http.StatusInternalServerError
	//	}
	//	return nil
	//}
	//
	//proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
