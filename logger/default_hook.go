package logger

import (
	"github.com/sirupsen/logrus"

	"github.com/blackdreamers/core/config"
)

type defaultHook struct{}

func (d defaultHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (d defaultHook) Fire(e *logrus.Entry) error {
	e.Data["app"] = config.Service.Name
	return nil
}

func init() {
	AddHook(&defaultHook{})
}
