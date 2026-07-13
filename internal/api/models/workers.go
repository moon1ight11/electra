package models

import "github.com/google/uuid"

// создание работника
type CreateWorkerRequest struct {
	Name           string `json:"name" binding:"required"`
	Phone          string `json:"phone" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Specialization string `json:"specialization"`
}

// запрос на логин
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// статистика работника
type WorkerStats struct {
	WorkerID       uuid.UUID `json:"worker_id"`
	WorkerName     string    `json:"worker_name"`
	OrdersCount    int       `json:"orders_count"`
	TotalEarned    float64   `json:"total_earned"`
	TotalTimeSpent int       `json:"total_time_spent"`
}

// общая статистика
type SummaryStats struct {
	OrdersCount    int     `json:"orders_count"`
	TotalEarned    float64 `json:"total_earned"`
	TotalTimeSpent int     `json:"total_time_spent"`
}
