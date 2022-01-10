package config

import (
	"strings"

	"github.com/blackdreamers/core/consts"
	"github.com/blackdreamers/core/env"
)

var (
	Log *logConf
)

type logConf struct {
	Level string
}

func (l *logConf) init() error {
	// etcd service配置的日志等级权重高于env中配置的
	l.Level = Service.Get(strings.ToLower(consts.LogLevel)).String(l.Level)
	return nil
}

func init() {
	Log = &logConf{
		Level: env.GetString(consts.LogLevel, "info"),
	}
	Configs(Log)
}
