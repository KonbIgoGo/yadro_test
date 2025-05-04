package main

import (
	"biathlon/config"
	"biathlon/internal/app"

	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()

	if err != nil {
		log.Fatalf("cannot get application config: %s", err)
	}

	var logger *zap.Logger
	logger, err = zap.NewProduction()

	if err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}

	err = app.Run(logger, cfg)
	if err != nil {
		log.Fatalf("processing stage error: %s", err)
	}
}
