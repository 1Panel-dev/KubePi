package server

import (
	"embed"
	"fmt"
	"github.com/KubeOperator/ekko/internal/config"
	v1Config "github.com/KubeOperator/ekko/internal/model/v1/config"
	"github.com/KubeOperator/ekko/migrate"
	"github.com/KubeOperator/ekko/pkg/file"
	"github.com/KubeOperator/ekko/pkg/i18n"
	"github.com/asdine/storm/v3"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/view"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"strings"
)

const sessionCookieName = "SESS_COOKIE_EKKO"

var EmbedWebEkko embed.FS
var EmbedWebDashboard embed.FS
var EmbedWebTerminal embed.FS

type EkkoSerer struct {
	*iris.Application
	db     *storm.DB
	logger *logrus.Logger
	config v1Config.Config
}

func NewEkkoSerer() *EkkoSerer {
	c := &EkkoSerer{}
	return c.bootstrap()
}

func (e *EkkoSerer) setUpConfig() {
	c, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	e.config = c
}

func (e *EkkoSerer) setUpLogger() {
	e.logger = logrus.New()
	l, err := logrus.ParseLevel(e.config.Spec.Logger.Level)
	if err != nil {
		e.logger.Errorf("cant not parse logger level %s, %s,use default level: INFO", e.config.Spec.Logger.Level, err)
	}
	e.logger.SetLevel(l)
}

func (e *EkkoSerer) setUpDB() {
	realDir := file.ReplaceHomeDir(e.config.Spec.DB.Path)
	if !fileutil.Exist(realDir) {
		if err := os.MkdirAll(realDir, 0755); err != nil {
			panic(fmt.Errorf("can not create database dir: %s message: %s", e.config.Spec.DB.Path, err))
		}
	}
	d, err := storm.Open(path.Join(realDir, "ekko.db"))
	if err != nil {
		panic(err)
	}
	e.db = d
}

func (e *EkkoSerer) makeEkkoDir() {
	ekkoDir := "~/.ekko"
	if !fileutil.Exist("~/.ekko") {
		if err := os.MkdirAll(ekkoDir, 0755); err != nil {
			panic(err)
		}
	}
}

func (e *EkkoSerer) setUpStaticFile() {
	e.Get("/", func(ctx *context.Context) {
		ctx.Redirect("/ekko")
	})
	spaOption := iris.DirOptions{SPA: true, IndexName: "index.html"}
	dashboardFS := iris.PrefixDir("web/dashboard", http.FS(EmbedWebDashboard))
	e.RegisterView(view.HTML(dashboardFS, ".html"))
	e.HandleDir("/dashboard/", dashboardFS, spaOption)
	terminalFS := iris.PrefixDir("web/terminal", http.FS(EmbedWebTerminal))
	e.RegisterView(view.HTML(terminalFS, ".html"))
	e.HandleDir("/terminal/", terminalFS, spaOption)
	ekkoFS := iris.PrefixDir("web/ekko", http.FS(EmbedWebEkko))
	e.RegisterView(view.HTML(ekkoFS, ".html"))
	e.HandleDir("/ekko", ekkoFS, spaOption)

}

func (e *EkkoSerer) setUpSession() {
	sess := sessions.New(sessions.Config{Cookie: sessionCookieName, AllowReclaim: true})
	e.Use(sess.Handler())
}

func (e *EkkoSerer) setResultHandler() {
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

func (e *EkkoSerer) setUpErrHandler() {
	e.OnAnyErrorCode(func(ctx iris.Context) {
		if ctx.Values().GetString("message") == "" {
			switch ctx.GetStatusCode() {
			case iris.StatusNotFound:
				ctx.Values().Set("message", "the server could not find the requested resource")
			}
		}
		message := ctx.Values().GetString("message")
		lang := ctx.Values().GetString("language")
		m, err := i18n.Translate(lang, message)
		if err != nil {
			e.Logger().Error(err)
			m = message
		}
		er := iris.Map{
			"success": false,
			"code":    ctx.GetStatusCode(),
			"message": m,
		}
		_, _ = ctx.JSON(er)
	})
}

func (e *EkkoSerer) runMigrations() {
	migrate.RunMigrate(e.db, e.logger)
}

func (e *EkkoSerer) bootstrap() *EkkoSerer {
	e.Application = iris.New()
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

var es *EkkoSerer

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
	es = NewEkkoSerer()
	route(es.Application)
	return es.Run(iris.Addr(fmt.Sprintf("%s:%d", es.config.Spec.Server.Bind.Host, es.config.Spec.Server.Bind.Port)))
}
