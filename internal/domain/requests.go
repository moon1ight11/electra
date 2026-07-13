package domain

import (
	"time"

	"github.com/google/uuid"
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
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Phone     string        `json:"phone"`
	Comment   *string       `json:"comment,omitempty"`
	Status    RequestStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}
