package interfaces

import (
	"context"
	"electra/internal/domain"
	"github.com/google/uuid"
)

// аутентификация и управление пользователями
type AuthService interface {
	// проверяет телефон и пароль, возвращает токен
	Login(ctx context.Context, phone, password string) (string, error)
	// создаёт нового исполнителя (только владелец)
	CreateWorker(ctx context.Context, ownerID uuid.UUID, name, phone, password string) (*domain.Worker, error)
}
