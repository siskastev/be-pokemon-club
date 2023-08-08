package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Out = os.Stdout
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func Logger() *logrus.Logger {
	return log
}
