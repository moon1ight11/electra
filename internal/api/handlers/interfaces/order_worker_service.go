package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"
	"github.com/google/uuid"
)

// отчётность исполнителей по заказам
type OrderWorkerService interface {
	// обновляет отчёт исполнителя по заказу
	UpdateReport(ctx context.Context, workerID uuid.UUID, input models.UpdateReportInput) error
	// возвращает отчёты всех исполнителей по заказу
	GetByOrder(ctx context.Context, userID, orderID uuid.UUID) ([]domain.OrderWorker, error)
	// снимает исполнителя с заказа
	RemoveWorker(ctx context.Context, ownerID, orderID, workerID uuid.UUID) error
	// возвращает историю выполненных заказов исполнителя
	ListCompletedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error)
	// возвращает все выполненные заказы (владелец)
	ListAllCompleted(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error)
}
