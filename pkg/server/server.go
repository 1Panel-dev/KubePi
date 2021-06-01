package server

import (
	"fmt"
	"github.com/KubeOperator/ekko/pkg/config"
	v1Config "github.com/KubeOperator/ekko/pkg/model/v1/config"
	"github.com/KubeOperator/ekko/pkg/route"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/sirupsen/logrus"
)

const sessionCookieName = "SESS_COOKIE_EKKO"

type EkkoSerer struct {
	*iris.Application
	db     *storm.DB
	logger *logrus.Logger
	config v1Config.Config
}

func (e *EkkoSerer) DB() *storm.DB {
	return e.db
}

func (e *EkkoSerer) Config() v1Config.Config {
	return e.config
}

func (e *EkkoSerer) Logger() *logrus.Logger {
	return e.logger
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
	route.InitRoute(e.Application)
	return e
}

func (e *EkkoSerer) Listen() error {
	return e.Run(iris.Addr(fmt.Sprintf("%s:%d", e.config.Spec.Server.Bind.Host, e.config.Spec.Server.Bind.Port)))
}
