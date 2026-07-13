package interfaces

import (
	"context"
	"electra/internal/domain"
	"github.com/google/uuid"
)

type AuthService interface {
	// проверка телефона и пароля и возврат токена
	Login(ctx context.Context, phone, password string) (string, error)
	// создание нового исполнителя
	CreateWorker(ctx context.Context, ownerID uuid.UUID, name, phone, password, specialization string) (*domain.Worker, error)
	// получение списка работников
	ListWorkers(ctx context.Context) ([]domain.WorkerInfo, error)
	// получение информации о пользователе
	GetMe(ctx context.Context, userID uuid.UUID) (*domain.WorkerInfo, error)
}
