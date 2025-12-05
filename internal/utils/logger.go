package utils

import (
	"log"

	"go.uber.org/zap"
)


var Logger *zap.Logger

func InitLogger(){
	var err error

	Logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger")
		panic(err)
	}

	defer Logger.Sync()
}