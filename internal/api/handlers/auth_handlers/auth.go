package authhandlers

import (
	"electra/internal/api/handlers/interfaces"
	"electra/pkg/logger"
)

type AuthHandler struct {
	authService interfaces.AuthService
	logger      logger.Logger
}

func NewAuthHandler(authService interfaces.AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}
