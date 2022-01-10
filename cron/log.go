package cron

import (
	"sync"

	"github.com/blackdreamers/core/consts"
	"github.com/blackdreamers/core/logger"
)

var (
	log  *cronLog
	once sync.Once
)

type cronLog struct{}

func newLog() *cronLog {
	once.Do(func() {
		log = &cronLog{}
	})

	return log
}

func (c cronLog) Info(msg string, keysAndValues ...interface{}) {
	logger.Fields(keysAndValues...).Log(logger.InfoLevel, msg)
}

func (c cronLog) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, consts.ErrKey, err)
	logger.Fields(keysAndValues...).Log(logger.ErrorLevel, msg)
}
