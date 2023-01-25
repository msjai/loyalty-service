package logger

import (
	"go.uber.org/zap"
)

func New() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // очищение всех буферезированных записей журналаа
	sugar := logger.Sugar()

	return sugar
}
