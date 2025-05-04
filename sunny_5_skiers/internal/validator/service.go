package validator

import (
	"biathlon/config"
	"biathlon/internal/processor"

	"go.uber.org/zap"
)

type implementation struct {
	logger    *zap.Logger
	cfg       *config.Config
	processor processor.Processor
}

func New(logger *zap.Logger, cfg *config.Config, processor processor.Processor) *implementation {
	return &implementation{
		logger:    logger,
		cfg:       cfg,
		processor: processor,
	}
}
