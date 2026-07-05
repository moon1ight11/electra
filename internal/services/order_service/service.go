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

// создать заказ
func (s *OrderService) CreateDirect(ctx context.Context, ownerID uuid.UUID, input models.CreateOrderDirectInput) (*domain.Order, error) {
	var plannedDate *time.Time
	if input.PlannedDate != "" {
		t, err := time.Parse("2006-01-02", input.PlannedDate)
		if err != nil {
			return nil, fmt.Errorf("error in create direct: invalid planned_date format: %w", err)
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
		return nil, fmt.Errorf("error in create direct: %w", err)
	}

	return order, nil
}

// запланированные заказы работника
func (s *OrderService) ListPlannedByWorker(ctx context.Context, workerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderRepo.ListPlannedByWorker(ctx, workerID)
	if err != nil {
		return nil, fmt.Errorf("error in list planned by worker: %w", err)
	}

	return orders, nil
}

// все запланированные заказы
func (s *OrderService) ListAllPlanned(ctx context.Context, ownerID uuid.UUID) ([]domain.Order, error) {
	orders, err := s.orderRepo.ListAllPlanned(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in list all planned: %w", err)
	}

	return orders, nil
}

// пометка выполнено
func (s *OrderService) Complete(ctx context.Context, workerID, orderID uuid.UUID) error {
	// проверяем, что исполнитель назначен на этот заказ
	ow, err := s.orderWorkerRepo.GetByOrderAndWorker(ctx, orderID, workerID)
	if err != nil {
		return fmt.Errorf("error in complete order (check assignment): %w", err)
	}
	if ow == nil {
		return errors.New("error in complete order: worker not assigned to this order")
	}

	// проверяем, что заказ ещё не завершён
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("error in complete order: %w", err)
	}
	if order == nil {
		return errors.New("error in complete order:: order not found")
	}
	if order.CompletedAt != nil {
		return errors.New("error in complete order:: order already completed")
	}

	err = s.orderRepo.Complete(ctx, orderID)
	if err != nil {
		return fmt.Errorf("error in complete order: %w", err)
	}

	return nil
}
