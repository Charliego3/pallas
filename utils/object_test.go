package utils

import (
	"github.com/charliego3/logger"
	"testing"
)

func TestNils(t *testing.T) {
	log := logger.Default().(*logger.DefaultLog)
	log.SetPrefix("[Application]")
	log.SetReportCaller(true)
	log.SetLevel(logger.LevelDebug)
	log.Debug("debug message", "key", "value")
	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")
	log.Fatal("fatal message")
}
