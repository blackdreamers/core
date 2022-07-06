package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

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

func Caller(skip int) func(f *runtime.Frame) (string, string) {
	return func(f *runtime.Frame) (string, string) {
		_, file, line, ok := runtime.Caller(skip)
		fileline := "unknown"
		if ok {
			filePath := strings.ReplaceAll(file, fmt.Sprintf("%s/pkg/mod/", os.Getenv("GOPATH")), "")
			// 去除路径中版本号
			versionIndex := strings.Index(filePath, "@")
			if versionIndex != -1 {
				subPath := filePath[versionIndex:]
				version := subPath[:strings.Index(subPath, "/")]
				filePath = strings.ReplaceAll(filePath, version, "")
			}
			fileline = fmt.Sprintf("%v:%v", filePath, line)
			if config.IsDevEnv() {
				fileline += "\t"
			}
		}
		return "", fileline
	}
}
