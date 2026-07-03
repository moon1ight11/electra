package orderworkerservice

import (
	"context"
	"errors"
	"fmt"

	"electra/internal/api/models"
	"electra/internal/domain"

	"github.com/google/uuid"
)

func (s *OrderWorkerService) UpdateReport(ctx context.Context, workerID uuid.UUID, input models.UpdateReportInput) error {
	ow, err := s.orderWorkerRepo.GetByOrderAndWorker(ctx, input.OrderID, workerID)
	if err != nil {
		return fmt.Errorf("get order worker: %w", err)
	}
	if ow == nil {
		return errors.New("worker not assigned to this order")
	}

	ow.TimeSpent = input.TimeSpent
	ow.EarnedAmount = input.EarnedAmount
	ow.MaterialsUsed = input.MaterialsUsed
	ow.Notes = input.Notes

	return s.orderWorkerRepo.Update(ctx, ow)
}

func (s *OrderWorkerService) GetByOrder(ctx context.Context, userID, orderID uuid.UUID) ([]domain.OrderWorker, error) {
	// работник видит только если он назначен на заказ
	ow, err := s.orderWorkerRepo.GetByOrderAndWorker(ctx, orderID, userID)
	if err != nil {
		return nil, fmt.Errorf("check access: %w", err)
	}
	if ow == nil {
		return nil, errors.New("access denied: not assigned to this order")
	}

	return s.orderWorkerRepo.GetByOrder(ctx, orderID)
}

func (s *OrderWorkerService) RemoveWorker(ctx context.Context, ownerID, orderID, workerID uuid.UUID) error {
	return s.orderWorkerRepo.Remove(ctx, orderID, workerID)
}

func (s *OrderWorkerService) ListCompletedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error) {
	return s.orderWorkerRepo.ListCompletedByWorker(ctx, workerID)
}

func (s *OrderWorkerService) ListAllCompleted(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error) {
	return s.orderWorkerRepo.ListAllCompleted(ctx)
}
