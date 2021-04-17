package logger

type Logger interface {
	init() error
}

func Init() error {
	if err := _logrus.init(); err != nil {
		return err
	}
	return nil
}
