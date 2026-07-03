package authservice

import (
	"electra/internal/api/jwt"
	"electra/internal/storage/repo/workers"
)

type AuthService struct {
	workerRepo *workers.WorkerRepo
	jwtService jwt.TokenService
}

func NewAuthService(workerRepo *workers.WorkerRepo, jwtService jwt.TokenService) *AuthService {
	return &AuthService{
		workerRepo: workerRepo,
		jwtService: jwtService,
	}
}
