package app

import (
	"biathlon/config"
	"biathlon/internal/processor"
	"biathlon/internal/validator"
	"bufio"
	"os"

	"go.uber.org/zap"
)

func Run(logger *zap.Logger, cfg *config.Config) error {

	events, err := os.Open("events")
	if err != nil {
		logger.Error("cannot open events file", zap.Error(err))
		return err
	}
	defer events.Close()

	validator := validator.New(logger, cfg, processor.New(cfg, logger))
	reader := bufio.NewReader(events)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		err = validator.Validate(line)
		if err != nil {
			logger.Error("failed to validate event", zap.Error(err))
		}
	}

	validator.GetLog()
	validator.GetResult()

	return nil
}
