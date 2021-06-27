package config

import (
	"strings"

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
	// etcd service配置的日志等级权重高于env中配置的
	l.Level = Service.Get(strings.ToLower(constant.LogLevel)).String(l.Level)
	return nil
}

func init() {
	Log = &logConf{
		Level: env.GetString(constant.LogLevel, logger.InfoLevel.String()),
	}
	Configs(Log)
}
