package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"

	"github.com/google/uuid"
)

type RequestService interface {
	Create(ctx context.Context, name, phone, comment string) (*domain.Request, error)
	ListNew(ctx context.Context, ownerID uuid.UUID) ([]domain.Request, error)
	ListAll(ctx context.Context, ownerID uuid.UUID) ([]domain.Request, error)
	ConvertToOrder(ctx context.Context, ownerID uuid.UUID, input models.CreateOrderFromRequestInput) (*domain.Order, error)
	Cancel(ctx context.Context, ownerID uuid.UUID, requestID uuid.UUID) error
}
