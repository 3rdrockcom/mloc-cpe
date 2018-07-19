package log

import (
	"github.com/epointpayment/mloc-cpe/app/config"

	"github.com/sirupsen/logrus"
)

var environment = config.EnvDevelopment

func init() {
	SetMode(environment)
}

func Start() {}

func Stop() {}

func SetMode(env string) {
	switch env {
	case config.EnvProduction:
		environment = env
		setProduction()
	case config.EnvDevelopment:
		environment = env
		setDevelopment()
	}
}

func setDevelopment() {
	DefaultLogger.SetLevel(logrus.DebugLevel)
	DefaultLogger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
}

func setProduction() {
	DefaultLogger.SetLevel(logrus.InfoLevel)
	DefaultLogger.Formatter = &logrus.JSONFormatter{}
}
