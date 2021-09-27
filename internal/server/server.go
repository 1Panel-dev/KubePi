package server

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

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

var EmbedWebKubePi embed.FS
var EmbedWebDashboard embed.FS
var EmbedWebTerminal embed.FS

type KubePiSerer struct {
	*iris.Application
	db     *storm.DB
	logger *logrus.Logger
	config v1Config.Config
}

func NewKubePiSerer() *KubePiSerer {
	c := &KubePiSerer{}
	return c.bootstrap()
}

func (e *KubePiSerer) setUpConfig() {
	c, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	e.config = c
}

func (e *KubePiSerer) setUpLogger() {
	e.logger = logrus.New()
	l, err := logrus.ParseLevel(e.config.Spec.Logger.Level)
	if err != nil {
		e.logger.Errorf("cant not parse logger level %s, %s,use default level: INFO", e.config.Spec.Logger.Level, err)
	}
	e.logger.SetLevel(l)
}

func (e *KubePiSerer) setUpDB() {
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

func (e *KubePiSerer) setUpStaticFile() {
	e.Get("/kubepi/", func(ctx *context.Context) {
		ctx.Redirect("/kubepi/kubepi")
	})
	spaOption := iris.DirOptions{SPA: true, IndexName: "index.html"}

	dashboardFS := iris.PrefixDir("web/dashboard", http.FS(EmbedWebDashboard))
	e.RegisterView(view.HTML(dashboardFS, ".html"))
	e.HandleDir("/kubepi/dashboard/", dashboardFS, spaOption)

	terminalFS := iris.PrefixDir("web/terminal", http.FS(EmbedWebTerminal))
	e.RegisterView(view.HTML(terminalFS, ".html"))
	e.HandleDir("/kubepi/terminal/", terminalFS, spaOption)

	kubePiFS := iris.PrefixDir("web/kubepi", http.FS(EmbedWebKubePi))
	e.RegisterView(view.HTML(kubePiFS, ".html"))
	e.HandleDir("/kubepi/kubepi/", kubePiFS, spaOption)
}

func (e *KubePiSerer) setUpSession() {
	sess := sessions.New(sessions.Config{Cookie: sessionCookieName, AllowReclaim: true})
	e.Use(sess.Handler())
}

func (e *KubePiSerer) setResultHandler() {
	e.Use(func(ctx *context.Context) {
		ctx.Next()
		isProxyPath := func() bool {
			p := ctx.GetCurrentRoute().Path()
			ss := strings.Split(p, "/")
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
				resp := iris.Map{
					"success": true,
					"data":    ctx.Values().Get("data"),
				}
				_, _ = ctx.JSON(resp)
			}
		}
	})
}

func (e *KubePiSerer) setUpErrHandler() {
	e.OnAnyErrorCode(func(ctx iris.Context) {
		if ctx.Values().GetString("message") == "" {
			switch ctx.GetStatusCode() {
			case iris.StatusNotFound:
				ctx.Values().Set("message", "the server could not find the requested resource")
			}
		}
		message := ctx.Values().Get("message")
		lang := ctx.Values().GetString("language")
		var (
			translateMessage string
			err              error
			originMessage    string
		)

		switch message.(type) {
		case string:
			originMessage = message.(string)
			translateMessage, err = i18n.Translate(lang, message.(string))
		case []string:
			ms := message.([]string)
			originMessage = strings.Join(ms, ",")
			if len(ms) > 0 {
				translateMessage, err = i18n.Translate(lang, ms[0], ms[1:])
			}
		}
		msg := translateMessage
		if err != nil {
			e.Logger().Warn(err)
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

func (e *KubePiSerer) runMigrations() {
	migrate.RunMigrate(e.db, e.logger)
}

func (e *KubePiSerer) bootstrap() *KubePiSerer {
	e.Application = iris.New()
	e.Application.Use(iris.Compression)
	e.setUpStaticFile()
	e.setUpConfig()
	e.setUpLogger()
	e.setUpDB()
	e.setUpSession()
	e.setResultHandler()
	e.setUpErrHandler()
	e.runMigrations()
	return e
}

var es *KubePiSerer

func DB() *storm.DB {
	return es.db
}

func Config() v1Config.Config {
	return es.config
}

func Logger() *logrus.Logger {
	return es.logger
}

func Listen(route func(party iris.Party)) error {
	es = NewKubePiSerer()
	route(es.Application)
	return es.Run(iris.Addr(fmt.Sprintf("%s:%d", es.config.Spec.Server.Bind.Host, es.config.Spec.Server.Bind.Port)))
}
