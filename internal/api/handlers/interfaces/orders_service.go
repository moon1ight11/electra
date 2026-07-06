package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"
	"github.com/google/uuid"
)

// работа с заказами
type OrderService interface {
	// создаёт заказ без заявки (по звонку)
	CreateDirect(ctx context.Context, ownerID uuid.UUID, input models.CreateOrderDirectInput) (*domain.Order, error)
	// возвращает запланированные заказы исполнителя
	ListPlannedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error)
	// возвращает все запланированные заказы (владелец)
	ListAllPlanned(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error)
	// завершает заказ
	Complete(ctx context.Context, workerID uuid.UUID, orderID uuid.UUID) error
	CompleteByOwner(ctx context.Context, orderID uuid.UUID) error
}
