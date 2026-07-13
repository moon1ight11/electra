package orderworkerservice

import (
	"context"
	"errors"
	"fmt"

	"electra/internal/api/models"
	"electra/internal/domain"

	"github.com/google/uuid"
)

// обновление отчета
func (s *OrderWorkerService) UpdateReport(ctx context.Context, workerID uuid.UUID, input models.UpdateReportInput) error {
	ow, err := s.orderWorkerRepo.GetByOrderAndWorker(ctx, input.OrderID, workerID)
	if err != nil {
		return fmt.Errorf("error in update report: %w", err)
	}
	if ow == nil {
		return errors.New("error in update report: worker not assigned to this order")
	}

	ow.TimeSpent = input.TimeSpent
	ow.EarnedAmount = input.EarnedAmount
	ow.MaterialsUsed = input.MaterialsUsed
	ow.Notes = input.Notes

	err = s.orderWorkerRepo.Update(ctx, ow)
	if err != nil {
		return fmt.Errorf("error in update report: %w", err)
	}

	return nil
}

// получение отчета из заказа
func (s *OrderWorkerService) GetByOrder(ctx context.Context, userID, orderID uuid.UUID) ([]domain.OrderWorker, error) {
	ow, _ := s.orderWorkerRepo.GetByOrderAndWorker(ctx, orderID, userID)
	if ow != nil {
		return s.orderWorkerRepo.GetByOrder(ctx, orderID)
	}

	worker, err := s.workerRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error in get user: %w", err)
	}
	if worker != nil && worker.Role == domain.RoleOwner {
		return s.orderWorkerRepo.GetByOrder(ctx, orderID)
	}

	return nil, errors.New("access denied: not assigned to this order")
}

// получение списка выполненных работ исполнителя
func (s *OrderWorkerService) ListCompletedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderWorkerRepo.ListCompletedByWorker(ctx, workerID)
	if err != nil {
		return nil, fmt.Errorf("error in list completed by worker: %w", err)
	}

	return orders, nil
}

// получение списка всех выполненных работ
func (s *OrderWorkerService) ListAllCompleted(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderWorkerRepo.ListAllCompleted(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in list all completed: %w", err)
	}

	return orders, nil
}

// удаление работника из заказа
func (s *OrderWorkerService) RemoveWorker(ctx context.Context, ownerID, orderID, workerID uuid.UUID) error {
	err := s.orderWorkerRepo.Remove(ctx, orderID, workerID)
	if err != nil {
		return fmt.Errorf("error in remove worker: %w", err)
	}

	return nil
}
