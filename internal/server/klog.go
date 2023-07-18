package server

import "github.com/go-logr/logr"

// 屏蔽 klog 日志
type TodoLogger struct {
	sink  logr.LogSink
	level int
}

func (t *TodoLogger) setSink(sink logr.LogSink) {
}

func (t TodoLogger) GetSink() logr.LogSink {
	return t.sink
}

func (t TodoLogger) WithSink(sink logr.LogSink) logr.Logger {
	return logr.Logger{}
}

func (t TodoLogger) Enabled() bool {
	return true
}

func (t TodoLogger) Info(msg string, keysAndValues ...interface{}) {

}

func (t TodoLogger) Error(err error, msg string, keysAndValues ...interface{}) {
}

func (t TodoLogger) V(level int) logr.Logger {
	return logr.Logger{}
}

func (t TodoLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	return logr.Logger{}
}

func (t TodoLogger) WithName(name string) logr.Logger {
	return logr.Logger{}
}

func (t TodoLogger) WithCallDepth(depth int) logr.Logger {
	return logr.Logger{}
}

func (t TodoLogger) WithCallStackHelper() (func(), logr.Logger) {
	return nil, logr.Logger{}
}
