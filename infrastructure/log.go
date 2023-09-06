package infrastructure

import (
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"os"
)

var Log = initializedLog()

func initializedLog() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(1)
	}
	log.Out = file

	log.SetFormatter(&ecslogrus.Formatter{})
	return log
}
