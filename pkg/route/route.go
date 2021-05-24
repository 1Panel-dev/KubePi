package route

import (
	"embed"
	"github.com/KubeOperator/ekko/pkg/route/proxy"
	"github.com/kataras/iris/v12"
	"net/http"
)

var EmbedWeb embed.FS

func InitRoute(app *iris.Application) {
	app.Party("/proxy").Any("/{p:path}", proxy.KubernetesApiProxy)
	fsys := iris.PrefixDir("web", http.FS(EmbedWeb))
	app.RegisterView(iris.HTML(fsys, ".html"))
	app.HandleDir("/ui", fsys)
}
