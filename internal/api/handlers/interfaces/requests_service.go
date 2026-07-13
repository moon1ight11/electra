package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"

	"github.com/google/uuid"
)

type RequestService interface {
	// создание заявки
	Create(ctx context.Context, name, phone, comment string) (*domain.Request, error)
	// получение новых заявок
	ListNew(ctx context.Context, ownerID uuid.UUID) ([]domain.Request, error)
	// получение всех заявок
	ListAll(ctx context.Context, ownerID uuid.UUID) ([]domain.Request, error)
	// пометка заявки как обработанной
	ConvertToOrder(ctx context.Context, ownerID uuid.UUID, input models.CreateOrderFromRequestInput) (*domain.Order, error)
	// пометка заявки как отмененной
	Cancel(ctx context.Context, ownerID uuid.UUID, requestID uuid.UUID) error
}
