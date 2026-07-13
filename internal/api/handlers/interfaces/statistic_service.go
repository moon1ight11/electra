package interfaces

import (
	"context"
	"electra/internal/api/models"
	"github.com/google/uuid"
)

type StatisticsService interface {
	// статистика по работнику за период
	ByWorker(ctx context.Context, ownerID, workerID uuid.UUID, from, to string) (*models.WorkerStats, error)
	// общая статистика по всем работникам за период
	AllWorkers(ctx context.Context, ownerID uuid.UUID, from, to string) ([]models.WorkerStats, error)
	// общая статистика по всей компании
	Summary(ctx context.Context, from, to string) (*models.SummaryStats, error)
}
