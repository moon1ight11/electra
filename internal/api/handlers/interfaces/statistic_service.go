package interfaces

import (
	"context"
	"electra/internal/api/models"
	"github.com/google/uuid"
)

type StatisticsService interface {
	ByWorker(ctx context.Context, ownerID, workerID uuid.UUID, from, to string) (*models.WorkerStats, error)
	AllWorkers(ctx context.Context, ownerID uuid.UUID, from, to string) ([]models.WorkerStats, error)
	Summary(ctx context.Context, from, to string) (*models.SummaryStats, error)
}
