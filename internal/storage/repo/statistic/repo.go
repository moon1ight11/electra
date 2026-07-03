package statistic

import "electra/internal/storage/repo/database"

type StatisticsRepo struct {
	db *database.DataBase
}

func NewStatisticsRepo(db *database.DataBase) *StatisticsRepo {
	return &StatisticsRepo{db: db}
}
