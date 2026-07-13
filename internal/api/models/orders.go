package models

import "github.com/google/uuid"

// создание заявки
type CreateRequestInput struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Comment string `json:"comment"`
}

// заказ из заявки
type CreateOrderFromRequestInput struct {
	RequestID      uuid.UUID   `json:"request_id" binding:"required"`
	Address        string      `json:"address" binding:"required"`
	Description    string      `json:"description"`
	EstimatedPrice *float64    `json:"estimated_price"`
	PlannedDate    string      `json:"planned_date"`
	WorkerIDs      []uuid.UUID `json:"worker_ids" binding:"required"`
}

// заказ мимо заявки
type CreateOrderDirectInput struct {
	Address        string      `json:"address" binding:"required"`
	Description    string      `json:"description"`
	EstimatedPrice *float64    `json:"estimated_price"`
	PlannedDate    string      `json:"planned_date"`
	WorkerIDs      []uuid.UUID `json:"worker_ids" binding:"required"`
}

// обновление отчета
type UpdateReportInput struct {
	OrderID       uuid.UUID `json:"order_id" binding:"required"`
	TimeSpent     *int      `json:"time_spent"`
	EarnedAmount  *float64  `json:"earned_amount"`
	MaterialsUsed *string   `json:"materials_used"`
	Notes         *string   `json:"notes"`
}

// обновление полей заказа
type UpdateOrderInput struct {
	ID             uuid.UUID `json:"id" binding:"required"`
	Address        string    `json:"address" binding:"required"`
	Description    string    `json:"description"`
	EstimatedPrice *float64  `json:"estimated_price"`
	PlannedDate    string    `json:"planned_date"`
}
