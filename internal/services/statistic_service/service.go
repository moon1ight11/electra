package statisticservice

import (
	"context"
	"fmt"
	"time"

	"electra/internal/api/models"

	"github.com/google/uuid"
)

// статистика по работнику
func (s *StatisticsService) ByWorker(ctx context.Context, ownerID, workerID uuid.UUID, from, to string) (*models.WorkerStats, error) {
	// парсим время периода
	fromTime, toTime, err := parsePeriod(from, to)
	if err != nil {
		return nil, fmt.Errorf("error in get stats by worker: %w", err)
	}

	row, err := s.statisticsRepo.ByWorker(ctx, workerID, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("error in get stats by worker: %w", err)
	}

	return &models.WorkerStats{
		WorkerID:       row.WorkerID,
		WorkerName:     row.WorkerName,
		OrdersCount:    row.OrdersCount,
		TotalEarned:    row.TotalEarned,
		TotalTimeSpent: row.TotalTimeSpent,
	}, nil
}

// статистика по всем работникам
func (s *StatisticsService) AllWorkers(ctx context.Context, ownerID uuid.UUID, from, to string) ([]models.WorkerStats, error) {
	// парсим время периода
	fromTime, toTime, err := parsePeriod(from, to)
	if err != nil {
		return nil, fmt.Errorf("error in get stats by all: %w", err)
	}

	rows, err := s.statisticsRepo.AllWorkers(ctx, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("error in get stats by all: %w", err)
	}

	result := make([]models.WorkerStats, len(rows))
	for i, r := range rows {
		result[i] = models.WorkerStats{
			WorkerID:       r.WorkerID,
			WorkerName:     r.WorkerName,
			OrdersCount:    r.OrdersCount,
			TotalEarned:    r.TotalEarned,
			TotalTimeSpent: r.TotalTimeSpent,
		}
	}

	return result, nil
}

// хэлпер для парсинга времени
func parsePeriod(from, to string) (time.Time, time.Time, error) {
	// парсинг времени от
	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		fromTime = time.Time{}
	}

	// парсинг времени до, если не вышло - то сутки
	toTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		toTime = time.Now().Add(24 * time.Hour)
	}

	return fromTime, toTime, nil
}
