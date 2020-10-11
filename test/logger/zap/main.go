package main

import (
	"github.com/macheal/go-micro-plugins/logger/zap"
	"github.com/micro/go-micro/v2/logger"
	uzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {

	NewLogger()
	for {
		logger.Info(logger.InfoLevel, "info")
		time.Sleep(10 * time.Millisecond)
	}

}

func NewLogger() {

	encoderConfig := uzap.NewProductionEncoderConfig()
	// 时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	conf := uzap.NewProductionConfig()
	conf.EncoderConfig = encoderConfig

	//转换日志等级
	loggerLevel, err := logger.GetLevel("debug")
	if err != nil {
		loggerLevel = logger.WarnLevel
	}
	l, err := zap.NewLogger(
		zap.WithConfig(conf),
		logger.WithLevel(loggerLevel),
	)
	if err != nil {
		panic(err)
	}
	logger.DefaultLogger = l

}
