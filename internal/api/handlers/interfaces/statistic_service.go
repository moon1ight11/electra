package interfaces

import (
	"context"
	"electra/internal/api/models"
	"github.com/google/uuid"
)

// статистика для владельца
type StatisticsService interface {
	// возвращает статистику по исполнителю за период
	ByWorker(ctx context.Context, ownerID, workerID uuid.UUID, from, to string) (*models.WorkerStats, error)
	// возвращает общую статистику по всем исполнителям за период
	AllWorkers(ctx context.Context, ownerID uuid.UUID, from, to string) ([]models.WorkerStats, error)
	// общая статистика по конторе
	Summary(ctx context.Context, from, to string) (*models.SummaryStats, error)
}
