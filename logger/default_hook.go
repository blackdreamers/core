package logger

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/blackdreamers/core/config"
)

type defaultHook struct{}

func (d defaultHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (d defaultHook) Fire(e *logrus.Entry) error {
	pos := "unknown"
	_, file, line, ok := runtime.Caller(6)
	if ok {
		pos = fmt.Sprintf("%v:%v", filepath.Base(file), line)
	}
	e.Data["pos"] = pos
	e.Data["app"] = config.Service.Name

	return nil
}

func init() {
	AddHook(&defaultHook{})
}
