package logger

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func InitLogger() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	logFile := &lumberjack.Logger{
		Filename:   "./logs/auth.log",
		MaxAge:     7,
		MaxBackups: 5,
		MaxSize:    100,
		Compress:   true,
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetLevel(logrus.InfoLevel)
}

func GetLogger() *logrus.Logger {
	return log
}
