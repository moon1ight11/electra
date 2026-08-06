package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateDirect(ctx context.Context, ownerID uuid.UUID, input models.CreateOrderDirectInput) (*domain.Order, error)
	ListPlannedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error)
	ListAllPlanned(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error)
	Complete(ctx context.Context, workerID uuid.UUID, orderID uuid.UUID) error
	CompleteByOwner(ctx context.Context, orderID uuid.UUID) error
	Update(ctx context.Context, input models.UpdateOrderInput) (*domain.Order, error)
}
