package jwt

import (
	"electra/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId    *uuid.UUID        `json:"user_id"`
	UserName  string            `json:"name"`
	UserPhone string            `json:"phone"`
	Role      domain.WorkerRole `json:"role"`
	jwt.RegisteredClaims
}
