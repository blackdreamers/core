package logger

import (
	"github.com/blackdreamers/go-micro/v3/logger"
)

func Init() error {
	var err error
	if logger.DefaultLogger, err = newLogrus(); err != nil {
		return err
	}
	return nil
}
