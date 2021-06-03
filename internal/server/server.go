package server

import (
	"embed"
	"errors"
	"fmt"
	"github.com/KubeOperator/ekko/internal/config"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Config "github.com/KubeOperator/ekko/internal/model/v1/config"
	"github.com/KubeOperator/ekko/internal/model/v1/user"
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/view"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
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

func (e *EkkoSerer) setResultHandler() {
	e.Use(func(ctx *context.Context) {
		ctx.Next()
		if ctx.GetStatusCode() >= iris.StatusOK && ctx.GetStatusCode() < iris.StatusBadRequest {
			resp := iris.Map{
				"success": true,
				"data":    ctx.Values().Get("data"),
			}
			_, _ = ctx.JSON(resp)
		}
	})
}

func (e *EkkoSerer) setSuperUser() {
	var superUser user.User
	if err := e.db.One("Name", "admin", &superUser); err != nil {
		if !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Sprintf("can not query supper user please check db connection: %s", err.Error()))
		}
	}

	if superUser.Name == "" {

		e.logger.Info("creat supper user")
		superUser = user.User{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "User",
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
			},
			Metadata: v1.Metadata{
				Name: "admin",
				UUID: uuid.New().String(),
			},
			Spec: user.Spec{
				Info: user.Info{
					NickName: "administrator",
					Email:    "support@fit2cloud.com",
				},
			},
		}
		pass := "admin123"
		hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost) //加密处理
		superUser.Spec.Authenticate.Password = string(hash)
		if err := e.db.Save(&superUser); err != nil {
			panic("can not save supper user please check db connection")
		}
		e.logger.Info("create supper user success")
		e.logger.Info("username: admin")
		e.logger.Info("password: admin123")
	}
}

func (e *EkkoSerer) setUpErrHandler() {
	e.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"success": false,
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
	e.setResultHandler()
	e.setUpErrHandler()
	e.setSuperUser()
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
