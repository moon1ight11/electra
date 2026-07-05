package statisticshandlers

import (
	"electra/internal/api/handlers/interfaces"
	"electra/pkg/logger"
)

type StatisticsHandler struct {
	statisticsService interfaces.StatisticsService
	logger            logger.Logger
}

func NewStatisticsHandler(statisticsService interfaces.StatisticsService, logger logger.Logger) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
		logger:            logger,
	}
}
