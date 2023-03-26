package log

import "go.uber.org/zap"

type Interface interface {
	Error(err error)
}

type logger struct {
	log *zap.Logger
}

func Init() Interface {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return &logger{log: zapLogger}
}

func (l *logger) Error(err error) {
	l.log.Error(err.Error())
}
