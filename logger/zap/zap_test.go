package zap

import (
	"testing"

	"github.com/macheal/go-micro/v2/logger"
)

func TestName(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	if l.String() != "zap" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("example logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	logger.DefaultLogger = l
	logger.Logf(logger.InfoLevel, "example logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	logger.DefaultLogger = l

	logger.Init(logger.WithLevel(logger.DebugLevel))
	l.Logf(logger.DebugLevel, "example show debug: %s", "debug msg")

	logger.Init(logger.WithLevel(logger.InfoLevel))
	l.Logf(logger.DebugLevel, "example non-show debug: %s", "debug msg")
}
