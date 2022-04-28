package logging

import "go.uber.org/zap"

func SugarLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
		}
	}(logger)
	sugar := logger.Sugar()
	return sugar
}
