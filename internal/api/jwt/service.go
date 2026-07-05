package jwt

import (
	"electra/internal/domain"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"regexp"
	"time"
)

type Service struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTService(secret string, expiration time.Duration) TokenService {
	return &Service{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

// валидация кастомных полей клеймов
func (c *Claims) CustomFieldsValidate() error {
	// проверяем валидность uuid
	if c.UserId == nil {
		return fmt.Errorf("error in customFieldsValidate: user id is empty")
	}

	// проверяем, что имя пользователя не пустое
	if c.UserName == "" {
		return fmt.Errorf("error in customFieldsValidate: invalid user name")
	}

	// проверяем валидность номера телефона
	var nonDigits = regexp.MustCompile(`\D`)

	cleaned := nonDigits.ReplaceAllString(c.UserPhone, "")
	if len(cleaned) == 0 {
		return fmt.Errorf("error in customFieldsValidate: phone must contain at least one digit")
	}

	return nil
}

// создание токена
func (j *Service) GenerateToken(userID uuid.UUID, userName string, userPhone string, userRole domain.WorkerRole) (string, error) {
	claims := &Claims{
		UserId:    &userID,
		UserName:  userName,
		UserPhone: userPhone,
		Role:      userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
		},
	}

	// валидация кастомных полей
	if err := claims.CustomFieldsValidate(); err != nil {
		return "", fmt.Errorf("Error in GenerateToken: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// декодировка токена
func (j *Service) ParseToken(tokenString string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Error in Parse token: invalid method")
		}
		return j.secret, nil
	})
}
