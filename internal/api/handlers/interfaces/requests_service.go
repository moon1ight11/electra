package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"

	"github.com/google/uuid"
)

// работа с заявками от клиентов
type RequestService interface {
	// создаёт заявку с сайта (публично)
	Create(ctx context.Context, name, phone, comment string) (*domain.Request, error)
	// возвращает новые (необработанные) заявки
	ListNew(ctx context.Context, ownerID uuid.UUID) ([]domain.Request, error)
	// возвращает все заявки
	ListAll(ctx context.Context, ownerID uuid.UUID) ([]domain.Request, error)
	// помечает заявку как converted и создаёт заказ
	ConvertToOrder(ctx context.Context, ownerID uuid.UUID, input models.CreateOrderFromRequestInput) (*domain.Order, error)
	// помечает заявку как отменённую
	Cancel(ctx context.Context, ownerID uuid.UUID, requestID uuid.UUID) error
}