package server

import "github.com/go-logr/logr"

// 屏蔽 klog 日志
type TodoLogger struct {
}

func (t TodoLogger) Enabled() bool {
	return true
}

func (t TodoLogger) Info(msg string, keysAndValues ...interface{}) {

}

func (t TodoLogger) Error(err error, msg string, keysAndValues ...interface{}) {
}

func (t TodoLogger) V(level int) logr.Logger {
	return nil
}

func (t TodoLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	return nil
}

func (t TodoLogger) WithName(name string) logr.Logger {
	return nil
}
