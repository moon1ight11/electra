package authservice

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/domain"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var nonDigits = regexp.MustCompile(`\D`)

func (s *AuthService) CreateWorker(ctx context.Context, ownerID uuid.UUID, name, phone, password, specialization string) (*domain.Worker, error) {
	owner, err := s.workerRepo.GetByID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner: %w", err)
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
		return nil, fmt.Errorf("failed to hash password: %w", err)
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
		return nil, fmt.Errorf("failed to create worker: %w", err)
	}

	return worker, nil
}

func (s *AuthService) Login(ctx context.Context, phone, password string) (string, error) {
	cleaned := nonDigits.ReplaceAllString(phone, "")
	if len(cleaned) == 0 {
		return "", errors.New("phone must contain at least one digit")
	}

	worker, err := s.workerRepo.GetByPhone(ctx, cleaned)
	if err != nil {
		return "", fmt.Errorf("failed to login: %w", err)
	}
	if worker == nil {
		return "", errors.New("worker not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(worker.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	token, err := s.jwtService.GenerateToken(worker.ID, worker.Name, cleaned, worker.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) ListWorkers(ctx context.Context) ([]domain.WorkerInfo, error) {
	workers, err := s.workerRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list workers: %w", err)
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
			Role:           w.Role,
			Specialization: spec,
		}
	}
	return result, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*domain.WorkerInfo, error) {
	worker, err := s.workerRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}
	if worker == nil {
		return nil, errors.New("worker not found")
	}

	spec := ""
	if worker.Specialization != nil {
		spec = *worker.Specialization
	}

	return &domain.WorkerInfo{
		ID:             worker.ID,
		Name:           worker.Name,
		Role:           worker.Role,
		Specialization: spec,
	}, nil
}

func (s *AuthService) DeleteWorker(ctx context.Context, ownerID, workerID uuid.UUID) error {
	owner, err := s.workerRepo.GetByID(ctx, ownerID)
	if err != nil {
		return fmt.Errorf("failed to get owner: %w", err)
	}
	if owner == nil || owner.Role != domain.RoleOwner {
		return errors.New("only owner can delete workers")
	}

	worker, err := s.workerRepo.GetByID(ctx, workerID)
	if err != nil {
		return fmt.Errorf("failed to get worker: %w", err)
	}
	if worker == nil {
		return errors.New("worker not found")
	}
	if worker.ID == ownerID {
		return errors.New("owner cannot delete themselves")
	}

	return s.workerRepo.Delete(ctx, workerID)
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, input models.UpdateProfileInput) (*domain.WorkerInfo, error) {
	worker, err := s.workerRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}
	if worker == nil {
		return nil, errors.New("worker not found")
	}

	if input.Name != "" {
		worker.Name = input.Name
	}

	var spec *string
	if input.Specialization != "" {
		spec = &input.Specialization
	} else {
		spec = worker.Specialization
	}

	var hash *string
	if input.Password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashStr := string(h)
		hash = &hashStr
	}

	if err := s.workerRepo.Update(ctx, userID, worker.Name, spec, hash); err != nil {
		return nil, fmt.Errorf("failed to update worker: %w", err)
	}

	specStr := ""
	if spec != nil {
		specStr = *spec
	}

	return &domain.WorkerInfo{
		ID:             worker.ID,
		Name:           worker.Name,
		Role:           worker.Role,
		Specialization: specStr,
	}, nil
}