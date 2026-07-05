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
	// работник видит только если он назначен на заказ
	ow, err := s.orderWorkerRepo.GetByOrderAndWorker(ctx, orderID, userID)
	if err != nil {
		return nil, fmt.Errorf("error in get by order (check access): %w", err)
	}
	if ow == nil {
		return nil, errors.New("error in get by order (access denied): not assigned to this order")
	}

	orw, err := s.orderWorkerRepo.GetByOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error in get by worker: %w", err)
	}

	return orw, nil
}

// удаление работника из заказа
func (s *OrderWorkerService) RemoveWorker(ctx context.Context, ownerID, orderID, workerID uuid.UUID) error {
	err := s.orderWorkerRepo.Remove(ctx, orderID, workerID)
	if err != nil {
		return fmt.Errorf("error in remove worker: %w", err)
	}

	return nil
}

// список выполненного работником
func (s *OrderWorkerService) ListCompletedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderWorkerRepo.ListCompletedByWorker(ctx, workerID)
	if err != nil {
		return nil, fmt.Errorf("error in list completed by worker: %w", err)
	}

	return orders, nil
}

// список выполненного всего
func (s *OrderWorkerService) ListAllCompleted(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderWorkerRepo.ListAllCompleted(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in list all completed: %w", err)
	}

	return orders, nil
}
