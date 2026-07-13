package jwt

import (
	"electra/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	// генерация токена из клеймов
	GenerateToken(userId uuid.UUID, userName string, userPhone string, userRole domain.WorkerRole) (string, error)
	// парсинг токена в клеймы
	ParseToken(token string, claims *Claims) (*jwt.Token, error)
}
