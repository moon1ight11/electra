package requesthandlers

import (
	"electra/internal/api/handlers/interfaces"
	"electra/pkg/logger"
)

type RequestHandler struct {
	requestService interfaces.RequestService
	logger         logger.Logger
}

func NewRequestHandler(requestService interfaces.RequestService, logger logger.Logger) *RequestHandler {
	return &RequestHandler{
		requestService: requestService,
		logger:         logger,
	}
}
