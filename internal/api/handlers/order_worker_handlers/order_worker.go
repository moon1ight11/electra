package orderworkerhandlers

import (
	"electra/internal/api/handlers/interfaces"
	"electra/pkg/logger"
)

type OrderWorkerHandler struct {
	orderWorkerService interfaces.OrderWorkerService
	logger             logger.Logger
}

func NewOrderWorkerHandler(orderWorkerService interfaces.OrderWorkerService, logger logger.Logger) *OrderWorkerHandler {
	return &OrderWorkerHandler{
		orderWorkerService: orderWorkerService,
		logger:             logger,
	}
}
