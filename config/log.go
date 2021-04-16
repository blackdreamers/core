package config

import (
	"github.com/blackdreamers/core/constant"
	"github.com/blackdreamers/core/env"
	"github.com/blackdreamers/go-micro/v3/logger"
)

var (
	Log *logConf
)

type logConf struct {
	Level string
}

func (l *logConf) init() error {
	return nil
}

func init() {
	Log = &logConf{
		Level: env.GetString(constant.LogLevel, logger.InfoLevel.String()),
	}
	Configs(Log)
}
