package domain

import (
	"time"
	"github.com/google/uuid"
)

// роль работника
type WorkerRole string
const (
	RoleOwner  WorkerRole = "owner"
	RoleWorker WorkerRole = "worker"
)

// статус заявки
type RequestStatus string
const (
	RequestNew       RequestStatus = "new"
	RequestConverted RequestStatus = "converted"
	RequestCancelled RequestStatus = "cancelled"
)

// заявка
type Request struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Phone     string         `json:"phone"`
	Comment   *string        `json:"comment,omitempty"`
	Status    RequestStatus  `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
}

// работник
type Worker struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Phone        *string    `json:"phone,omitempty"`
	Role         WorkerRole `json:"role"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
}

// заказ
type Order struct {
	ID             uuid.UUID  `json:"id"`
	RequestID      *uuid.UUID `json:"request_id,omitempty"`
	Address        string     `json:"address"`
	Description    *string    `json:"description,omitempty"`
	EstimatedPrice *float64   `json:"estimated_price,omitempty"`
	PlannedDate    *time.Time `json:"planned_date,omitempty"`
	CreatedBy      uuid.UUID      `json:"created_by"`
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
