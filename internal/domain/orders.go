package domain

import (
	"time"

	"github.com/google/uuid"
)

// заказ
type Order struct {
	ID             uuid.UUID  `json:"id"`
	RequestID      *uuid.UUID `json:"request_id,omitempty"`
	RequestPhone   *string    `json:"request_phone,omitempty"`
	Address        string     `json:"address"`
	Description    *string    `json:"description,omitempty"`
	EstimatedPrice *float64   `json:"estimated_price,omitempty"`
	PlannedDate    *time.Time `json:"planned_date,omitempty"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// заказы работников
type OrderWorker struct {
	ID            uuid.UUID `json:"id"`
	OrderID       uuid.UUID `json:"order_id"`
	WorkerID      uuid.UUID `json:"worker_id"`
	TimeSpent     *int      `json:"time_spent,omitempty"`
	EarnedAmount  *float64  `json:"earned_amount,omitempty"`
	MaterialsUsed *string   `json:"materials_used,omitempty"`
	Notes         *string   `json:"notes,omitempty"`
}
