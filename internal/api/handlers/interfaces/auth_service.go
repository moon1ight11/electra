package interfaces

import (
	"context"
	"electra/internal/domain"
	"github.com/google/uuid"
)

type WorkerInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// аутентификация и управление пользователями
type AuthService interface {
	// проверяет телефон и пароль, возвращает токен
	Login(ctx context.Context, phone, password string) (string, error)
	// создаёт нового исполнителя (только владелец)
	CreateWorker(ctx context.Context, ownerID uuid.UUID, name, phone, password string) (*domain.Worker, error)
	ListWorkers(ctx context.Context) ([]WorkerInfo, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*WorkerInfo, error)
}
