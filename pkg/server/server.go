package server

import (
	"embed"
	"fmt"
	"github.com/KubeOperator/ekko/pkg/config"
	v1Config "github.com/KubeOperator/ekko/pkg/model/v1/config"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/view"
	"github.com/sirupsen/logrus"
	"net/http"
)

const sessionCookieName = "SESS_COOKIE_EKKO"

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
	d, err := storm.Open(e.config.Spec.DB.Path)
	if err != nil {
		panic(err)
	}
	e.db = d
}

func (e *EkkoSerer) setUpStaticFile() {
	dashboardFS := iris.PrefixDir("web/dashboard", http.FS(EmbedWebDashboard))
	e.RegisterView(view.HTML(dashboardFS, ".html"))
	e.HandleDir("/ui", dashboardFS)

	terminalFS := iris.PrefixDir("web/terminal", http.FS(EmbedWebTerminal))
	e.RegisterView(view.HTML(terminalFS, ".html"))
	e.HandleDir("/terminal", terminalFS)

}

func (e *EkkoSerer) setUpSession() {
	sess := sessions.New(sessions.Config{Cookie: sessionCookieName, AllowReclaim: true})
	e.Use(sess.Handler())
}

func (e *EkkoSerer) setUpErrHandler() {
	e.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}
		_, _ = ctx.JSON(err)
	})
}

func (e *EkkoSerer) bootstrap() *EkkoSerer {
	e.Application = iris.New()
	e.setUpConfig()
	e.setUpLogger()
	e.setUpDB()
	e.setUpSession()
	e.setUpErrHandler()
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
