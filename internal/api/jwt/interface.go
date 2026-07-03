package jwt

import (
	"electra/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(userId uuid.UUID, userName string, userPhone string, userRole domain.WorkerRole) (string, error)
	ParseToken(token string, claims *Claims) (*jwt.Token, error)
}
