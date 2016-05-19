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

type logrusLogger struct {
	level logrus.Level
}

func (l *logrusLogger) Printf(format string, v ...interface{}) {
	switch l.level {
	default:
		logrus.Printf(format, v...)
	case logrus.ErrorLevel:
		logrus.Errorf(format, v...)
	case logrus.WarnLevel:
		logrus.Warnf(format, v...)
	case logrus.InfoLevel:
		logrus.Infof(format, v...)
	case logrus.DebugLevel:
		logrus.Debugf(format, v...)
	}

}

var LogrusInfoLogger = &logrusLogger{
	level: logrus.InfoLevel,
}
var LogrusErrorLogger = &logrusLogger{
	level: logrus.ErrorLevel,
}
