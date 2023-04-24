package logger

import (
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-extras/elogrus.v7"

	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/consts/timef"
	"github.com/blackdreamers/core/utils"
)

var (
	log = &Logrus{hooks: make(logrus.LevelHooks)}
)

type Logrus struct {
	hooks logrus.LevelHooks
	entry *logrus.Entry
}

func (l *Logrus) init() error {
	level, err := logrus.ParseLevel(config.Log.Level)
	if err != nil {
		return err
	}

	std := logrus.New()

	if config.IsDevEnv() {
		std.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: timef.YearMonthDayHourMinuteSecond,
		})
	} else {
		std.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: timef.YearMonthDayHourMinuteSecond,
		})

		client, err := elasticsearch.NewClient(
			elasticsearch.Config{
				Addresses: []string{config.ES.Host},
				Username:  config.ES.User,
				Password:  config.ES.Password,
			},
		)
		if err != nil {
			return err
		}

		ip, err := utils.GetLocalIpV4()
		if err != nil {
			return err
		}

		hook, err := elogrus.NewAsyncElasticHook(client, ip, level, config.Log.Index)
		if err != nil {
			return err
		}
		AddHook(hook)
	}

	std.SetOutput(os.Stdout)
	std.SetLevel(level)
	std.ReplaceHooks(log.hooks)

	l.entry = logrus.NewEntry(std)

	return nil
}

func GetEntry() *logrus.Entry {
	return log.entry
}

func AddHook(hook logrus.Hook) {
	for _, level := range hook.Levels() {
		log.hooks[level] = append(log.hooks[level], hook)
	}
}
