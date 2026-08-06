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
		return fmt.Errorf("failed to get order worker: %w", err)
	}
	if ow == nil {
		return errors.New("worker not assigned to this order")
	}

	ow.TimeSpent = input.TimeSpent
	ow.EarnedAmount = input.EarnedAmount
	ow.MaterialsUsed = input.MaterialsUsed
	ow.Notes = input.Notes

	if err := s.orderWorkerRepo.Update(ctx, ow); err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}

	return nil
}

func (s *OrderWorkerService) GetByOrder(ctx context.Context, userID, orderID uuid.UUID) ([]domain.OrderWorker, error) {
	ow, _ := s.orderWorkerRepo.GetByOrderAndWorker(ctx, orderID, userID)
	if ow != nil {
		return s.orderWorkerRepo.GetByOrder(ctx, orderID)
	}

	worker, err := s.workerRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if worker != nil && worker.Role == domain.RoleOwner {
		return s.orderWorkerRepo.GetByOrder(ctx, orderID)
	}

	return nil, errors.New("access denied: not assigned to this order")
}

func (s *OrderWorkerService) ListCompletedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderWorkerRepo.ListCompletedByWorker(ctx, workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list completed orders: %w", err)
	}

	return orders, nil
}

func (s *OrderWorkerService) ListAllCompleted(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderWorkerRepo.ListAllCompleted(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all completed orders: %w", err)
	}

	return orders, nil
}

func (s *OrderWorkerService) RemoveWorker(ctx context.Context, ownerID, orderID, workerID uuid.UUID) error {
	if err := s.orderWorkerRepo.Remove(ctx, orderID, workerID); err != nil {
		return fmt.Errorf("failed to remove worker: %w", err)
	}

	return nil
}