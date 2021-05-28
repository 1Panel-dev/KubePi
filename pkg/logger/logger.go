package logger

import (
	"fmt"
	"github.com/KubeOperator/ekko/pkg/config"
	"github.com/sirupsen/logrus"
)

var instance *logrus.Logger

const loggerLevelParseErr = "cant not parse logger level %s, %s"

func Init() {
	c := config.ReadConfig()
	instance = logrus.New()
	l, err := logrus.ParseLevel(c.Spec.Logger.Level)
	if err != nil {
		panic(fmt.Sprintf(loggerLevelParseErr, c.Spec.Logger.Level, err))
	}
	instance.SetLevel(l)
}

func GetLogger() *logrus.Logger {
	return instance
}
