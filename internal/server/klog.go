package server

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
