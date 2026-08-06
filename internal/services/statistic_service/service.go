package statisticservice

import (
	"context"
	"fmt"
	"time"

	"electra/internal/api/models"

	"github.com/google/uuid"
)

func (s *StatisticsService) ByWorker(ctx context.Context, ownerID, workerID uuid.UUID, from, to string) (*models.WorkerStats, error) {
	fromTime, toTime, err := parsePeriod(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to parse period: %w", err)
	}

	row, err := s.statisticsRepo.ByWorker(ctx, workerID, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats by worker: %w", err)
	}

	return &models.WorkerStats{
		WorkerID:       row.WorkerID,
		WorkerName:     row.WorkerName,
		OrdersCount:    row.OrdersCount,
		TotalEarned:    row.TotalEarned,
		TotalTimeSpent: row.TotalTimeSpent,
	}, nil
}

func (s *StatisticsService) AllWorkers(ctx context.Context, ownerID uuid.UUID, from, to string) ([]models.WorkerStats, error) {
	fromTime, toTime, err := parsePeriod(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to parse period: %w", err)
	}

	rows, err := s.statisticsRepo.AllWorkers(ctx, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats all workers: %w", err)
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

func (s *StatisticsService) Summary(ctx context.Context, from, to string) (*models.SummaryStats, error) {
	fromTime, toTime, err := parsePeriod(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to parse period: %w", err)
	}

	row, err := s.statisticsRepo.SummaryStats(ctx, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary stats: %w", err)
	}

	return &models.SummaryStats{
		OrdersCount:    row.OrdersCount,
		TotalEarned:    row.TotalEarned,
		TotalTimeSpent: row.TotalTimeSpent,
	}, nil
}

func parsePeriod(from, to string) (time.Time, time.Time, error) {
	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		fromTime = time.Time{}
	}

	toTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		toTime = time.Now().Add(24 * time.Hour)
	}

	return fromTime, toTime, nil
}