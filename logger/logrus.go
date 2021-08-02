package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/core/constant/timef"
	plslog "github.com/blackdreamers/go-micro/plugins/logger/logrus/v3"
	"github.com/blackdreamers/go-micro/v3/logger"
)

var (
	_logrus = &logrus{hooks: make(log.LevelHooks)}
)

type logrus struct {
	hooks              log.LevelHooks
	enableCallerEntry  *log.Entry
	disableCallerEntry *log.Entry
}

func (l *logrus) init() error {
	level, err := log.ParseLevel(config.Log.Level)
	if err != nil {
		return err
	}

	enableCallerStd := newLogrus(level, true)
	disableCallerStd := newLogrus(level, false)

	l.enableCallerEntry = log.NewEntry(enableCallerStd)
	l.disableCallerEntry = log.NewEntry(disableCallerStd)

	logger.DefaultLogger = plslog.NewLogger(plslog.WithLogger(GetEntry(true)))

	return nil
}

func GetEntry(reportCaller bool) *log.Entry {
	if reportCaller {
		return _logrus.enableCallerEntry
	}
	return _logrus.disableCallerEntry
}

func newLogrus(level log.Level, reportCaller bool) *log.Logger {
	std := log.New()

	std.SetFormatter(&log.JSONFormatter{
		CallerPrettyfier: caller(8),
		TimestampFormat:  timef.YearMonthDayHourMinuteSecond,
	})
	if config.IsDevEnv() {
		std.SetFormatter(&log.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			CallerPrettyfier: caller(9),
			TimestampFormat:  timef.YearMonthDayHourMinuteSecond,
		})
	}

	std.SetOutput(os.Stdout)
	std.SetLevel(level)
	std.SetReportCaller(reportCaller)
	std.ReplaceHooks(_logrus.hooks)

	return std
}

func addHook(level log.Level, hks ...log.Hook) {
	_logrus.hooks[level] = hks
}

func caller(skip int) func(f *runtime.Frame) (string, string) {
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
