package main

import (
	"fmt"
	"github.com/macheal/go-micro-plugins/logger/zap"
	"github.com/micro/go-micro/v2/logger"
	"github.com/natefinch/lumberjack"
	uzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {

	NewLogger()
	for {
		logger.Log(logger.InfoLevel, "info")
		logger.Log(logger.DebugLevel, "debug")
		logger.Log(logger.ErrorLevel, "error")
		time.Sleep(10 * time.Millisecond)
	}

}

func NewLogger() {

	encoderConfig := uzap.NewProductionEncoderConfig()
	// 时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	conf := uzap.NewProductionConfig()
	conf.EncoderConfig = encoderConfig
	///dev/stderr
	conf.OutputPaths = []string{"stderr", "a.log"}
	fileName := "micro-srv.log"
	outer := lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   1, //128, //MB
		MaxAge:    7,
		LocalTime: true,
		Compress:  true,
	}
	fmt.Println(outer)

	//转换日志等级
	loggerLevel, err := logger.GetLevel("warn")
	if err != nil {
		loggerLevel = logger.WarnLevel
	}
	//loggerLevel = logger.WarnLevel
	l, err := zap.NewLogger(
		zap.WithConfig(conf),
		logger.WithLevel(loggerLevel),
		zap.WithOutput(&outer),
		zap.WithSingleOutputOutput(true),
	)
	if err != nil {
		panic(err)
	}
	logger.DefaultLogger = l

}
