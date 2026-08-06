package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, phone, password string) (string, error)
	CreateWorker(ctx context.Context, ownerID uuid.UUID, name, phone, password, specialization string) (*domain.Worker, error)
	ListWorkers(ctx context.Context) ([]domain.WorkerInfo, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*domain.WorkerInfo, error)
	DeleteWorker(ctx context.Context, ownerID, workerID uuid.UUID) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) (*domain.WorkerInfo, error)
}
