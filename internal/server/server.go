package server

import (
	"embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"

	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
	"k8s.io/klog/v2"

	"github.com/KubeOperator/kubepi/internal/config"
	v1Config "github.com/KubeOperator/kubepi/internal/model/v1/config"
	"github.com/KubeOperator/kubepi/migrate"
	"github.com/KubeOperator/kubepi/pkg/file"
	"github.com/KubeOperator/kubepi/pkg/i18n"
	"github.com/asdine/storm/v3"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/view"
	"github.com/sirupsen/logrus"
)

const sessionCookieName = "SESS_COOKIE_KUBEPI"

var SessionMgr *sessions.Sessions

var EmbedWebKubePi embed.FS
var EmbedWebDashboard embed.FS
var EmbedWebTerminal embed.FS
var WebkubectlEntrypoint string

type Option func(server *KubePiServer)

func WithServerBindHost(host string) Option {
	return func(server *KubePiServer) {
		if host != "" {
			server.config.Spec.Server.Bind.Host = host
		}
	}
}

func WithServerBindPort(port int) Option {
	return func(server *KubePiServer) {
		if port != 0 {
			server.config.Spec.Server.Bind.Port = port
		}
	}
}

func WithCustomConfigFilePath(path string) Option {
	return func(server *KubePiServer) {
		if path != "" {
			server.configCustomFilePath = path
		}
	}
}

type KubePiServer struct {
	app                  *iris.Application
	db                   *storm.DB
	logger               *logrus.Logger
	configCustomFilePath string
	config               *v1Config.Config
	rootRoute            iris.Party
}

func NewKubePiSerer(opts ...Option) *KubePiServer {
	c := &KubePiServer{}
	c.app = iris.New()
	c.config = getDefaultConfig()
	for _, op := range opts {
		op(c)
	}
	c.setUpConfig()
	for _, op := range opts {
		op(c)
	}
	return c.bootstrap()
}

func (e *KubePiServer) setUpConfig() {
	err := config.ReadConfig(e.config, e.configCustomFilePath)
	if err != nil {
		panic(err)
	}
}

func (e *KubePiServer) setUpLogger() {
	klog.SetLogger(TodoLogger{})
	e.logger = logrus.New()
	l, err := logrus.ParseLevel(e.config.Spec.Logger.Level)
	if err != nil {
		e.logger.Errorf("cant not parse logger level %s, %s,use default level: INFO", e.config.Spec.Logger.Level, err)
	}
	e.logger.SetLevel(l)
}

func (e *KubePiServer) setUpDB() {
	realDir := file.ReplaceHomeDir(e.config.Spec.DB.Path)
	if !fileutil.Exist(realDir) {
		if err := os.MkdirAll(realDir, 0755); err != nil {
			panic(fmt.Errorf("can not create database dir: %s message: %s", e.config.Spec.DB.Path, err))
		}
	}
	d, err := storm.Open(path.Join(realDir, "kubepi.db"))
	if err != nil {
		panic(err)
	}
	e.db = d
}

func (e *KubePiServer) setUpRootRoute() {
	e.app.Any("/", func(ctx *context.Context) {
		ctx.Redirect("/kubepi")
	})
	c := swagger.Config{
		URL: "/kubepi/swagger/doc.json",
	}
	e.app.Get("/kubepi/swagger/{any:path}", swagger.CustomWrapHandler(&c, swaggerFiles.Handler))
	e.rootRoute = e.app.Party("/kubepi")
}

func (e *KubePiServer) setUpStaticFile() {
	spaOption := iris.DirOptions{SPA: true, IndexName: "index.html"}
	party := e.rootRoute.Party("/")
	party.Use(iris.Compression)
	dashboardFS := iris.PrefixDir("web/dashboard", http.FS(EmbedWebDashboard))
	party.RegisterView(view.HTML(dashboardFS, ".html"))
	party.HandleDir("/dashboard/", dashboardFS, spaOption)

	terminalFS := iris.PrefixDir("web/terminal", http.FS(EmbedWebTerminal))
	party.RegisterView(view.HTML(terminalFS, ".html"))
	party.HandleDir("/terminal/", terminalFS, spaOption)

	kubePiFS := iris.PrefixDir("web/kubepi", http.FS(EmbedWebKubePi))
	party.RegisterView(view.HTML(kubePiFS, ".html"))
	party.HandleDir("/", kubePiFS, spaOption)
}

func (e *KubePiServer) setUpSession() {
	SessionMgr = sessions.New(sessions.Config{Cookie: sessionCookieName, AllowReclaim: true, Expires: time.Duration(e.config.Spec.Session.Expires) * time.Hour})
	e.rootRoute.Use(SessionMgr.Handler())
}

const ContentTypeDownload = "application/download"

func (e *KubePiServer) setResultHandler() {
	e.rootRoute.Use(func(ctx *context.Context) {
		ctx.Next()
		contentType := ctx.ResponseWriter().Header().Get("Content-Type")
		if contentType == ContentTypeDownload {
			return
		}
		isProxyPath := func() bool {
			p := ctx.GetCurrentRoute().Path()
			ss := strings.Split(p, "/")
			if len(ss) > 0 {
				if ss[0] == "webkubectl" {
					return true
				}
			}
			if len(ss) >= 3 {
				for i := range ss {
					if ss[i] == "proxy" || ss[i] == "ws" {
						return true
					}
				}
			}
			return false
		}()
		if !isProxyPath {
			if ctx.GetStatusCode() >= iris.StatusOK && ctx.GetStatusCode() < iris.StatusBadRequest {
				if ctx.Values().Get("token") != nil {
					_, _ = ctx.Write(ctx.Values().Get("token").([]uint8))
				} else {
					resp := iris.Map{
						"success": true,
						"data":    ctx.Values().Get("data"),
					}
					_, _ = ctx.JSON(resp)
				}
			}
		}
	})
}

func (e *KubePiServer) setUpErrHandler() {
	e.rootRoute.OnAnyErrorCode(func(ctx iris.Context) {
		if ctx.Values().GetString("message") == "" {
			switch ctx.GetStatusCode() {
			case iris.StatusNotFound:
				ctx.Values().Set("message", "the server could not find the requested resource")
			}
		}
		message := ctx.Values().Get("message")
		if message == nil || message == "" {
			message = ctx.Values().Get("iris.context.error")
		}

		lang := ctx.Values().GetString("language")
		var (
			translateMessage string
			err              error
			originMessage    string
		)

		switch value := message.(type) {
		case string:
			originMessage = message.(string)
			translateMessage, err = i18n.Translate(lang, value)
		case []string:
			originMessage = strings.Join(value, ",")
			if len(value) > 0 {
				translateMessage, err = i18n.Translate(lang, value[0], value[1:])
			}
		case context.ErrPrivate:
			err := message.(context.ErrPrivate)
			translateMessage = err.Error()
		}
		msg := translateMessage
		if err != nil {
			e.app.Logger().Debug(err)
			msg = originMessage
		}
		er := iris.Map{
			"success": false,
			"code":    ctx.GetStatusCode(),
			"message": msg,
		}
		_, _ = ctx.JSON(er)
	})
}

func (e *KubePiServer) runMigrations() {
	migrate.RunMigrate(e.db, e.logger)
}
func (e *KubePiServer) setWebkubectlProxy() {
	handler := func(ctx *context.Context) {
		p := ctx.Params().Get("p")
		if strings.Contains(p, "root") {
			ctx.Request().URL.Path = strings.ReplaceAll(ctx.Request().URL.Path, "root", "")
			ctx.Request().RequestURI = strings.ReplaceAll(ctx.Request().RequestURI, "root", "")
		}
		u, _ := url.Parse("http://localhost:8080")
		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ModifyResponse = func(resp *http.Response) error {
			if resp.StatusCode == iris.StatusMovedPermanently {
				// 重定向重写
				if resp.Header.Get("Location") == "/kubepi/webkubectl/" {
					resp.Header.Set("Location", "/kubepi/webkubectl/root")
				}
			}
			return nil
		}
		proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	}
	e.rootRoute.Any("/webkubectl/{p:path}", handler)
	e.rootRoute.Any("webkubectl", handler)
}

func (e *KubePiServer) setUpTtyEntrypoint() {
	f, err := os.OpenFile("init-kube.sh", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		e.logger.Error(err)
		return
	}
	if _, err := f.WriteString(WebkubectlEntrypoint); err != nil {
		e.logger.Error(err)
		return
	}
}

func (e *KubePiServer) bootstrap() *KubePiServer {
	e.setUpRootRoute()
	e.setUpStaticFile()
	e.setUpLogger()
	e.setUpDB()
	e.setUpSession()
	e.setResultHandler()
	e.setUpErrHandler()
	e.setWebkubectlProxy()
	e.runMigrations()
	e.setUpTtyEntrypoint()
	e.startTty()
	return e
}

var es *KubePiServer

func DB() *storm.DB {
	return es.db
}

func Config() *v1Config.Config {
	return es.config
}

func Logger() *logrus.Logger {
	return es.logger
}

func Listen(route func(party iris.Party), options ...Option) error {
	es = NewKubePiSerer(options...)
	route(es.rootRoute)
	return es.app.Run(iris.Addr(fmt.Sprintf("%s:%d", es.config.Spec.Server.Bind.Host, es.config.Spec.Server.Bind.Port)))
}

func getDefaultConfig() *v1Config.Config {
	return &v1Config.Config{
		BaseModel: v1.BaseModel{
			ApiVersion: "v1",
			Kind:       "AppConfig",
		},
		Metadata: v1.Metadata{},
		Spec: v1Config.Spec{
			Server: v1Config.ServerConfig{
				Bind: v1Config.BindConfig{
					Host: "0.0.0.0",
					Port: 80,
				},
				SSL: v1Config.SSLConfig{
					Enable: false,
				},
			},
			DB: v1Config.DBConfig{
				Path: "/var/lib/kubepi/db",
			},
			Session: v1Config.SessionConfig{
				Expires: 72,
			},
			Logger: v1Config.LoggerConfig{Level: "debug"},
			Jwt:    v1Config.JwtConfig{},
		},
	}
}
