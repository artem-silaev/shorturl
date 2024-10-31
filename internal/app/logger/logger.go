package logger

import (
	"go.uber.org/zap"
)

var Log zap.SugaredLogger

func Init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	Log = *logger.Sugar()
}
