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

// работник
type Worker struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Phone          *string    `json:"phone,omitempty"`
	Role           WorkerRole `json:"role"`
	PasswordHash   string     `json:"-"`
	Specialization *string    `json:"specialization,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	DeletedAt      *time.Time `json:"-"`
}

// информация по работнику
type WorkerInfo struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Role           WorkerRole `json:"role"`
	Specialization string    `json:"specialization,omitempty"`
}

// суммарная статистика
type SummaryRow struct {
	OrdersCount    int
	TotalEarned    float64
	TotalTimeSpent int
}

// статичтика по работнику
type StatsRow struct {
	WorkerID       uuid.UUID
	WorkerName     string
	OrdersCount    int
	TotalEarned    float64
	TotalTimeSpent int
}
