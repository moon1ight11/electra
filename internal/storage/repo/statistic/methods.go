package statistic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StatsRow struct {
	WorkerID       uuid.UUID
	WorkerName     string
	OrdersCount    int
	TotalEarned    float64
	TotalTimeSpent int
}

type SummaryRow struct {
	OrdersCount    int
	TotalEarned    float64
	TotalTimeSpent int
}

// статистика по исполнителю за период
func (r *StatisticsRepo) ByWorker(ctx context.Context, workerID uuid.UUID, from, to time.Time) (*StatsRow, error) {
	query := `
		SELECT w.id, w.name,
		       COUNT(DISTINCT ow.order_id),
		       COALESCE(SUM(ow.earned_amount), 0),
		       COALESCE(SUM(ow.time_spent), 0)
		FROM order_workers ow
		JOIN orders o ON o.id = ow.order_id
		JOIN workers w ON w.id = ow.worker_id
		WHERE ow.worker_id = $1
		  AND o.completed_at >= $2
		  AND o.completed_at <= $3
		GROUP BY w.id, w.name`

	row := &StatsRow{}
	err := r.db.DB.QueryRowContext(ctx, query, workerID, from, to).Scan(
		&row.WorkerID,
		&row.WorkerName,
		&row.OrdersCount,
		&row.TotalEarned,
		&row.TotalTimeSpent)
	if err != nil {
		if err == sql.ErrNoRows {
			return &StatsRow{WorkerID: workerID}, nil
		}
		return nil, fmt.Errorf("error in stats by worker: %w", err)
	}
	return row, nil
}

// статистика по всем исполнителям за период
func (r *StatisticsRepo) AllWorkers(ctx context.Context, from, to time.Time) ([]StatsRow, error) {
	query := `
		SELECT w.id, w.name,
		       COUNT(DISTINCT ow.order_id),
		       COALESCE(SUM(ow.earned_amount), 0),
		       COALESCE(SUM(ow.time_spent), 0)
		FROM workers w
		LEFT JOIN order_workers ow ON ow.worker_id = w.id
		LEFT JOIN orders o ON o.id = ow.order_id
			AND o.completed_at >= $1
			AND o.completed_at <= $2
		GROUP BY w.id, w.name
		ORDER BY w.name`

	rows, err := r.db.DB.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("error in stats all workers: %w", err)
	}
	defer rows.Close()

	stats := make([]StatsRow, 0)
	for rows.Next() {
		var s StatsRow
		if err := rows.Scan(
			&s.WorkerID,
			&s.WorkerName,
			&s.OrdersCount,
			&s.TotalEarned,
			&s.TotalTimeSpent); err != nil {
			return nil, fmt.Errorf("error in scan stats: %w", err)
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

func (r *StatisticsRepo) SummaryStats(ctx context.Context, from, to time.Time) (*SummaryRow, error) {
	query := `
		SELECT
			COUNT(DISTINCT o.id),
			COALESCE(SUM(ow.earned_amount), 0),
			COALESCE(SUM(ow.time_spent), 0)
		FROM orders o
		LEFT JOIN order_workers ow ON ow.order_id = o.id
		WHERE o.completed_at >= $1 AND o.completed_at <= $2`

	row := &SummaryRow{}
	err := r.db.DB.QueryRowContext(ctx, query, from, to).
		Scan(&row.OrdersCount, &row.TotalEarned, &row.TotalTimeSpent)
	if err != nil {
		if err == sql.ErrNoRows {
			return &SummaryRow{}, nil
		}
		return nil, fmt.Errorf("summary stats: %w", err)
	}
	return row, nil
}