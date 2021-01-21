package log

import (
	"log"

	"github.com/masibw/blog-server/config"
	"go.uber.org/zap"
)

func GetPureLogger() (logger *zap.Logger) {
	var err error

	if config.IsLocal() {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return
}

func GetLogger() (sugar *zap.SugaredLogger) {
	logger := GetPureLogger()

	sugar = logger.Sugar()
	return sugar
}
