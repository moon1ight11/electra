package orderservice

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *OrderService) CreateDirect(ctx context.Context, ownerID uuid.UUID, input models.CreateOrderDirectInput) (*domain.Order, error) {
	var plannedDate *time.Time
	if input.PlannedDate != "" {
		t, err := time.Parse("2006-01-02", input.PlannedDate)
		if err != nil {
			return nil, fmt.Errorf("invalid planned_date format: %w", err)
		}
		plannedDate = &t
	}

	order := &domain.Order{
		Address:        input.Address,
		Description:    &input.Description,
		EstimatedPrice: input.EstimatedPrice,
		PlannedDate:    plannedDate,
		CreatedBy:      ownerID,
	}

	if err := s.orderRepo.Create(ctx, order, input.WorkerIDs); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *OrderService) ListPlannedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderRepo.ListPlannedByWorker(ctx, workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list planned orders: %w", err)
	}

	return orders, nil
}

func (s *OrderService) ListAllPlanned(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderRepo.ListAllPlanned(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all planned orders: %w", err)
	}

	return orders, nil
}

func (s *OrderService) Complete(ctx context.Context, workerID, orderID uuid.UUID) error {
	ow, err := s.orderWorkerRepo.GetByOrderAndWorker(ctx, orderID, workerID)
	if err != nil {
		return fmt.Errorf("failed to check assignment: %w", err)
	}
	if ow == nil {
		return errors.New("worker not assigned to this order")
	}

	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errors.New("order not found")
	}
	if order.CompletedAt != nil {
		return errors.New("order already completed")
	}

	return s.orderRepo.Complete(ctx, orderID)
}

func (s *OrderService) CompleteByOwner(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errors.New("order not found")
	}
	if order.CompletedAt != nil {
		return errors.New("order already completed")
	}

	return s.orderRepo.Complete(ctx, orderID)
}

func (s *OrderService) Update(ctx context.Context, input models.UpdateOrderInput) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.CompletedAt != nil {
		return nil, errors.New("cannot edit completed order")
	}

	var plannedDate *time.Time
	if input.PlannedDate != "" {
		t, err := time.Parse("2006-01-02", input.PlannedDate)
		if err != nil {
			return nil, fmt.Errorf("invalid planned_date format: %w", err)
		}
		plannedDate = &t
	}

	order.Address = input.Address
	order.Description = &input.Description
	order.EstimatedPrice = input.EstimatedPrice
	order.PlannedDate = plannedDate

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}