package workerhandlers

import "electra/internal/api/handlers/interfaces"

type WorkerHandler struct {
	authService interfaces.AuthService
}

func NewWorkerHandler(authService interfaces.AuthService) *WorkerHandler {
	return &WorkerHandler{authService: authService}
}