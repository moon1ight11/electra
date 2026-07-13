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

// создание работника
func (s *AuthService) CreateWorker(ctx context.Context, ownerID uuid.UUID, name, phone, password, specialization string) (*domain.Worker, error) {
	owner, err := s.workerRepo.GetByID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("error in get owner: %w", err)
	}
	if owner == nil || owner.Role != domain.RoleOwner {
		return nil, errors.New("only owner can create workers")
	}

	cleaned := nonDigits.ReplaceAllString(phone, "")
	if len(cleaned) == 0 {
		return nil, errors.New("phone must contain at least one digit")
	}

	existing, _ := s.workerRepo.GetByPhone(ctx, cleaned)
	if existing != nil {
		return nil, errors.New("worker with this phone already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error in hash password: %w", err)
	}

	var spec *string
	if specialization != "" {
		spec = &specialization
	}

	worker := &domain.Worker{
		Name:           name,
		Phone:          &cleaned,
		Role:           domain.RoleWorker,
		PasswordHash:   string(hash),
		Specialization: spec,
	}

	if err := s.workerRepo.Create(ctx, worker); err != nil {
		return nil, fmt.Errorf("error in create worker: %w", err)
	}

	return worker, nil
}

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

	if err := bcrypt.CompareHashAndPassword([]byte(worker.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	token, err := s.jwtService.GenerateToken(worker.ID, worker.Name, cleaned, worker.Role)
	if err != nil {
		return "", fmt.Errorf("error in generate token: %w", err)
	}

	return token, nil
}

// получение списка работников
func (s *AuthService) ListWorkers(ctx context.Context) ([]domain.WorkerInfo, error) {
	workers, err := s.workerRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in list workers: %w", err)
	}

	result := make([]domain.WorkerInfo, len(workers))
	for i, w := range workers {
		spec := ""
		if w.Specialization != nil {
			spec = *w.Specialization
		}
		result[i] = domain.WorkerInfo{
			ID:             w.ID,
			Name:           w.Name,
			Specialization: spec,
		}
	}
	return result, nil
}

// получение информации о текущем пользователе
func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*domain.WorkerInfo, error) {
	worker, err := s.workerRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error in get me: %w", err)
	}
	if worker == nil {
		return nil, errors.New("worker not found")
	}

	return &domain.WorkerInfo{
		ID:   worker.ID,
		Name: worker.Name,
	}, nil
}
