package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/KubeOperator/ekko/pkg/config"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	keyPrefix           = "Bearer"
	AuthorizationHeader = "Authorization"
)

func KubernetesApiProxy(ctx *context.Context) {
	c := config.GetConfig()
	proxyPath := ctx.Params().Get("p")
	u, err := url.Parse(c.KubeApiServerConfig.ApiServer)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	token := fmt.Sprintf("%s %s", keyPrefix, c.Token)
	fmt.Println(token)
	ctx.Request().Header.Add(AuthorizationHeader, token)
	ctx.Request().URL.Path = proxyPath
	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
