package statisticshandlers

import "electra/internal/api/handlers/interfaces"

type StatisticsHandler struct {
	statisticsService interfaces.StatisticsService
}

func NewStatisticsHandler(statisticsService interfaces.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statisticsService: statisticsService}
}