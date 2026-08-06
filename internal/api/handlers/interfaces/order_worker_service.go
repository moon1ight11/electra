package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"
	"github.com/google/uuid"
)

type OrderWorkerService interface {
	UpdateReport(ctx context.Context, workerID uuid.UUID, input models.UpdateReportInput) error
	GetByOrder(ctx context.Context, userID, orderID uuid.UUID) ([]domain.OrderWorker, error)
	RemoveWorker(ctx context.Context, ownerID, orderID, workerID uuid.UUID) error
	ListCompletedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error)
	ListAllCompleted(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error)
}
