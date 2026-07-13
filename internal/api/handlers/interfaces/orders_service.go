package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"
	"github.com/google/uuid"
)

type OrderService interface {
	// создание заказа без заявки
	CreateDirect(ctx context.Context, ownerID uuid.UUID, input models.CreateOrderDirectInput) (*domain.Order, error)
	// получение запланированных заказов исполнителя
	ListPlannedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error)
	// получение всех запланированных заказов
	ListAllPlanned(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error)
	// завершение заказа
	Complete(ctx context.Context, workerID uuid.UUID, orderID uuid.UUID) error
	// завершение заказа владельцем
	CompleteByOwner(ctx context.Context, orderID uuid.UUID) error
	// обновление заказа
	Update(ctx context.Context, input models.UpdateOrderInput) (*domain.Order, error)
}
