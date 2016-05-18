package common

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

func ConfigString(name string) string {
	return os.Getenv(fmt.Sprintf("BRC_%s", name))
}

func ConfigByteArray(name string) []byte {
	base64Value := ConfigString(name)
	value, err := base64.StdEncoding.DecodeString(base64Value)
	if err != nil {
		logrus.WithError(err).WithField("name", name).Warn("Error while decoding config value")
		return nil
	}
	return value
}
