package statisticservice

import "electra/internal/storage/repo/statistic"

type StatisticsService struct {
	statisticsRepo *statistic.StatisticsRepo
}

func NewStatisticsService(statisticsRepo *statistic.StatisticsRepo) *StatisticsService {
	return &StatisticsService{
		statisticsRepo: statisticsRepo,
	}
}
