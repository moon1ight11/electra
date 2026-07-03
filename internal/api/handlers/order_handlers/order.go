package orderhandlers

import "electra/internal/api/handlers/interfaces"

type OrderHandler struct {
	orderService interfaces.OrderService
}

func NewOrderHandler(orderService interfaces.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}
