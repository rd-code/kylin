package logger

import "testing"

func TestLoggerImpl_Info(t *testing.T) {
	Default.Level(Warn)
	Default.Info("=========")
	Default.Warn("=============%s", "xxx")
	Default.Error("===========---------,%s", "bbbb")
}
