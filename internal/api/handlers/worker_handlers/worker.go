package workerhandlers

import (
	"electra/internal/api/handlers/interfaces"
	"electra/pkg/logger"
)

type WorkerHandler struct {
	authService interfaces.AuthService
	logger      logger.Logger
}

func NewWorkerHandler(authService interfaces.AuthService, logger logger.Logger) *WorkerHandler {
	return &WorkerHandler{
		authService: authService,
		logger:      logger,
	}
}
