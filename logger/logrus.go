package logger

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/constant"
	log "github.com/blackdreamers/go-micro/plugins/logger/logrus/v3"
	"github.com/blackdreamers/go-micro/v3/logger"
)

var (
	hooks = make(logrus.LevelHooks)
	entry *logrus.Entry
)

func newLogrus() (logger.Logger, error) {
	std := logrus.New()
	std.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: constant.Timestamp,
	})
	if config.Conf.Env == constant.Prod {
		std.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: constant.Timestamp,
		})
	}
	std.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(config.Conf.LogLevel)
	if err != nil {
		return nil, err
	}
	std.SetLevel(level)
	std.ReplaceHooks(hooks)

	entry = logrus.NewEntry(std)

	return log.NewLogger(log.WithLogger(entry)), nil
}

func addHook(level logrus.Level, hks ...logrus.Hook) {
	hooks[level] = hks
}

func GetLogrus() *logrus.Entry {
	return entry
}
