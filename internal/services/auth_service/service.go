package authservice

import (
	"context"
	"electra/internal/domain"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var nonDigits = regexp.MustCompile(`\D`)

// логин
func (s *AuthService) Login(ctx context.Context, phone, password string) (string, error) {
	// чистим телефон от нецифровых символов
	cleaned := nonDigits.ReplaceAllString(phone, "")
	if len(cleaned) == 0 {
		return "", errors.New("phone must contain at least one digit")
	}

	worker, err := s.workerRepo.GetByPhone(ctx, cleaned)
	if err != nil {
		return "", fmt.Errorf("error in login: %w", err)
	}
	if worker == nil {
		return "", errors.New("error in login: worker not found")
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(worker.PasswordHash), []byte(password)); err != nil {
	// 	return "", errors.New("invalid password")
	// }

	if worker.PasswordHash != password {
		return "", errors.New("invalid password")
	}

	token, err := s.jwtService.GenerateToken(worker.ID, worker.Name, cleaned, worker.Role)
	if err != nil {
		return "", fmt.Errorf("error in generate token: %w", err)
	}

	return token, nil
}

// создание работника
func (s *AuthService) CreateWorker(ctx context.Context, ownerID uuid.UUID, name, phone, password string) (*domain.Worker, error) {
	// проверяем, что создаёт владелец
	owner, err := s.workerRepo.GetByID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("error in get owner: %w", err)
	}
	if owner == nil || owner.Role != domain.RoleOwner {
		return nil, errors.New("error in create worker: only owner can create workers")
	}

	// чистим телефон
	cleaned := nonDigits.ReplaceAllString(phone, "")
	if len(cleaned) == 0 {
		return nil, errors.New("phone must contain at least one digit")
	}

	// проверяем, нет ли уже такого телефона
	existing, _ := s.workerRepo.GetByPhone(ctx, cleaned)
	if existing != nil {
		return nil, errors.New("error: worker with this phone already exists")
	}

	// хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error in hash password: %w", err)
	}

	worker := &domain.Worker{
		Name:         name,
		Phone:        &cleaned,
		Role:         domain.RoleWorker,
		PasswordHash: string(hash),
	}

	if err := s.workerRepo.Create(ctx, worker); err != nil {
		return nil, fmt.Errorf("error in create worker: %w", err)
	}

	return worker, nil
}
