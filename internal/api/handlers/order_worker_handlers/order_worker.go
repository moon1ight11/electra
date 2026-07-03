package orderworkerhandlers

import "electra/internal/api/handlers/interfaces"

type OrderWorkerHandler struct {
	orderWorkerService interfaces.OrderWorkerService
}

func NewOrderWorkerHandler(orderWorkerService interfaces.OrderWorkerService) *OrderWorkerHandler {
	return &OrderWorkerHandler{orderWorkerService: orderWorkerService}
}