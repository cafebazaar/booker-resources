package common

import (
	"os"

	"github.com/Sirupsen/logrus"
)

func init() {
	logrus.SetOutput(os.Stderr)
	switch ConfigString("LOG_LEVEL") {
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	}
}
