package authservice

import (
	"context"
	"electra/internal/api/jwt"
	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/storage/repo/database"
	"electra/internal/storage/repo/workers"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupAuth(t *testing.T) (*AuthService, *workers.WorkerRepo, func()) {
	t.Helper()

	db, err := database.PostgresConnection(
		config.Config{
			DataBase: config.DataBaseConfig{
				Host:          "localhost",
				Port:          15432,
				DBName:        "electra",
				User:          "user",
				Password:      "user_pass",
				MigrationsDir: "../../storage/repo/migrations",
			},
		})
	if err != nil {
		t.Fatalf("db: %v", err)
	}

	db.UpMigrations()

	db.DB.Exec("DELETE FROM order_workers")
	db.DB.Exec("DELETE FROM orders")
	db.DB.Exec("DELETE FROM requests")
	db.DB.Exec("DELETE FROM workers")

	workerRepo := workers.NewWorkerRepo(db)
	jwtService := jwt.NewJWTService("test-secret", 24)
	authService := NewAuthService(workerRepo, jwtService)

	return authService, workerRepo, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func TestLogin_InvalidPhone(t *testing.T) {
	auth, _, cleanup := setupAuth(t)
	defer cleanup()

	_, err := auth.Login(context.Background(), "", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one digit")
}

func TestLogin_WorkerNotFound(t *testing.T) {
	auth, _, cleanup := setupAuth(t)
	defer cleanup()

	_, err := auth.Login(context.Background(), "79991112233", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestLogin_WrongPassword(t *testing.T) {
	auth, repo, cleanup := setupAuth(t)
	defer cleanup()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	repo.Create(context.Background(), &domain.Worker{
		Name: "Test", Phone: strPtr("79991112233"), Role: domain.RoleWorker, PasswordHash: string(hash),
	})

	_, err := auth.Login(context.Background(), "79991112233", "wrong")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")
}

func TestLogin_Success(t *testing.T) {
	auth, repo, cleanup := setupAuth(t)
	defer cleanup()

	hash, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.DefaultCost)
	repo.Create(context.Background(), &domain.Worker{
		Name: "Test", Phone: strPtr("79991112233"), Role: domain.RoleOwner, PasswordHash: string(hash),
	})

	token, err := auth.Login(context.Background(), "79991112233", "mypassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_PhoneWithSymbols(t *testing.T) {
	auth, repo, cleanup := setupAuth(t)
	defer cleanup()

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	repo.Create(context.Background(), &domain.Worker{
		Name: "Test", Phone: strPtr("79991112233"), Role: domain.RoleWorker, PasswordHash: string(hash),
	})

	token, err := auth.Login(context.Background(), "+7 (999) 111-22-33", "pass")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestCreateWorker_NotOwner(t *testing.T) {
	auth, repo, cleanup := setupAuth(t)
	defer cleanup()

	worker := &domain.Worker{Name: "Vasya", Phone: strPtr("79161112233"), Role: domain.RoleWorker, PasswordHash: "hash"}
	repo.Create(context.Background(), worker)

	_, err := auth.CreateWorker(context.Background(), worker.ID, "New", "79001234567", "pass", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only owner")
}

func TestCreateWorker_Success(t *testing.T) {
	auth, repo, cleanup := setupAuth(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	repo.Create(context.Background(), owner)

	w, err := auth.CreateWorker(context.Background(), owner.ID, "Petya", "79007654321", "pass", "Электрик")
	assert.NoError(t, err)
	assert.Equal(t, "Petya", w.Name)
	assert.Equal(t, "Электрик", *w.Specialization)
}

func TestCreateWorker_DuplicatePhone(t *testing.T) {
	auth, repo, cleanup := setupAuth(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	repo.Create(context.Background(), owner)

	auth.CreateWorker(context.Background(), owner.ID, "First", "79007654321", "pass", "")

	_, err := auth.CreateWorker(context.Background(), owner.ID, "Second", "79007654321", "pass", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func strPtr(s string) *string { return &s }
