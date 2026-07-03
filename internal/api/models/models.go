package models

import "github.com/google/uuid"

type CreateOrderFromRequestInput struct {
	RequestID      uuid.UUID   `json:"request_id" binding:"required"`
	Address        string      `json:"address" binding:"required"`
	Description    string      `json:"description"`
	EstimatedPrice *float64    `json:"estimated_price"`
	PlannedDate    string      `json:"planned_date"`
	WorkerIDs      []uuid.UUID `json:"worker_ids" binding:"required"`
}

type CreateOrderDirectInput struct {
	Address        string      `json:"address" binding:"required"`
	Description    string      `json:"description"`
	EstimatedPrice *float64    `json:"estimated_price"`
	PlannedDate    string      `json:"planned_date"`
	WorkerIDs      []uuid.UUID `json:"worker_ids" binding:"required"`
}

type UpdateReportInput struct {
	OrderID       uuid.UUID `json:"order_id" binding:"required"`
	TimeSpent     *int      `json:"time_spent"`
	EarnedAmount  *float64  `json:"earned_amount"`
	MaterialsUsed *string   `json:"materials_used"`
	Notes         *string   `json:"notes"`
}

type WorkerStats struct {
	WorkerID       uuid.UUID `json:"worker_id"`
	WorkerName     string    `json:"worker_name"`
	OrdersCount    int       `json:"orders_count"`
	TotalEarned    float64   `json:"total_earned"`
	TotalTimeSpent int       `json:"total_time_spent"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateWorkerRequest struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateRequestInput struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Comment string `json:"comment"`
}
