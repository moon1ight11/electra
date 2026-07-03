package requestservice

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *RequestService) Create(ctx context.Context, name, phone, comment string) (*domain.Request, error) {
	if len(phone) == 0 {
		return nil, errors.New("phone must contain at least one digit")
	}

	var c *string
	if comment != "" {
		c = &comment
	}

	req := &domain.Request{
		Name:    name,
		Phone:   phone,
		Comment: c,
	}

	if err := s.requestRepo.Create(ctx, req); err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return req, nil
}

func (s *RequestService) ListNew(ctx context.Context, ownerID uuid.UUID) ([]domain.Request, error) {
	return s.requestRepo.ListNew(ctx)
}

func (s *RequestService) ListAll(ctx context.Context, ownerID uuid.UUID) ([]domain.Request, error) {
	return s.requestRepo.List(ctx)
}

func (s *RequestService) ConvertToOrder(ctx context.Context, ownerID uuid.UUID, input models.CreateOrderFromRequestInput) (*domain.Order, error) {
	// проверяем, что заявка существует и новая
	req, err := s.requestRepo.GetByID(ctx, input.RequestID)
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}
	if req == nil {
		return nil, errors.New("request not found")
	}
	if req.Status != domain.RequestNew {
		return nil, errors.New("request already processed")
	}

	var plannedDate *time.Time
	if input.PlannedDate != "" {
		t, err := time.Parse("2006-01-02", input.PlannedDate)
		if err != nil {
			return nil, fmt.Errorf("invalid planned_date format: %w", err)
		}
		plannedDate = &t
	}

	order := &domain.Order{
		RequestID:      &input.RequestID,
		Address:        input.Address,
		Description:    &input.Description,
		EstimatedPrice: input.EstimatedPrice,
		PlannedDate:    plannedDate,
		CreatedBy:      ownerID,
	}

	if err := s.orderRepo.Create(ctx, order, input.WorkerIDs); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	// помечаем заявку как converted
	if err := s.requestRepo.MarkConverted(ctx, input.RequestID); err != nil {
		return nil, fmt.Errorf("mark converted: %w", err)
	}

	return order, nil
}

func (s *RequestService) Cancel(ctx context.Context, ownerID uuid.UUID, requestID uuid.UUID) error {
	return s.requestRepo.MarkCancelled(ctx, requestID)
}
