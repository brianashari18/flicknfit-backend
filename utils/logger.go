package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var appLogger *logrus.Logger

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}

func GetLogger() *logrus.Logger {
	if appLogger == nil {
		appLogger = NewLogger()
	}
	return appLogger
}
