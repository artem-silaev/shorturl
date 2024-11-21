package logger

import (
	"go.uber.org/zap"
)

var Log Logger

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
}

type ZapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return &ZapLogger{
		sugaredLogger: logger.Sugar(),
	}, nil
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

func Init() {
	logger, err := NewZapLogger()
	if err != nil {
		panic(err)
	}
	Log = logger
}
