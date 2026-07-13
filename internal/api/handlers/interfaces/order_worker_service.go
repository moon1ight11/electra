package interfaces

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"
	"github.com/google/uuid"
)

type OrderWorkerService interface {
	// обновление отчета исполнителя по заказу
	UpdateReport(ctx context.Context, workerID uuid.UUID, input models.UpdateReportInput) error
	// получение отчетов всех исполнителей по заказу
	GetByOrder(ctx context.Context, userID, orderID uuid.UUID) ([]domain.OrderWorker, error)
	// удаление исполнителя с заказа
	RemoveWorker(ctx context.Context, ownerID, orderID, workerID uuid.UUID) error
	// получение выполненных заказов исполнителем
	ListCompletedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error)
	// получение всех выполненных заказов
	ListAllCompleted(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error)
}
