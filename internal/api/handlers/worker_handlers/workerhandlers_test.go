package workerhandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"electra/internal/api/jwt"
	"electra/internal/api/models"
	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/services/auth_service"
	"electra/internal/storage/repo/database"
	"electra/internal/storage/repo/workers"
	"electra/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupWorkerHandler(t *testing.T) (*WorkerHandler, *workers.WorkerRepo, func()) {
	t.Helper()

	db, err := database.PostgresConnection(config.Config{
		DataBase: config.DataBaseConfig{
			Host:          "localhost",
			Port:          15432,
			DBName:        "electra",
			User:          "user",
			Password:      "user_pass",
			MigrationsDir: "../../../storage/repo/migrations",
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
	authService := authservice.NewAuthService(workerRepo, jwtService)
	log := logger.NewNopLogger()
	handler := NewWorkerHandler(authService, log)

	return handler, workerRepo, func() {
		db.Close()
	}
}

func TestListWorkers_Empty(t *testing.T) {
	handler, _, cleanup := setupWorkerHandler(t)
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/worker/workers", nil)

	handler.ListWorkers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}

func TestListWorkers_WithData(t *testing.T) {
	handler, repo, cleanup := setupWorkerHandler(t)
	defer cleanup()

	repo.Create(
		context.Background(),
		&domain.Worker{
			Name:         "Petya",
			Phone:        strPtr("79007654321"),
			Role:         domain.RoleWorker,
			PasswordHash: "hash",
		},
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/worker/workers", nil)

	handler.ListWorkers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Petya")
}

func TestGetMe(t *testing.T) {
	handler, repo, cleanup := setupWorkerHandler(t)
	defer cleanup()

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	wkr := &domain.Worker{
		Name:         "Petya",
		Phone:        strPtr("79007654321"),
		Role:         domain.RoleWorker,
		PasswordHash: string(hash),
	}
	repo.Create(context.Background(), wkr)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/worker/me", nil)
	c.Set("UserId", wkr.ID.String())

	handler.GetMe(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Petya")
}

func TestCreateWorker_ByOwner(t *testing.T) {
	handler, repo, cleanup := setupWorkerHandler(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	repo.Create(context.Background(), owner)

	body := models.CreateWorkerRequest{
		Name:     "New Worker",
		Phone:    "79009999999",
		Password: "pass123",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/owner/workers", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("UserId", owner.ID.String())

	handler.CreateWorker(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "New Worker")
}

func TestCreateWorker_NotOwner(t *testing.T) {
	handler, repo, cleanup := setupWorkerHandler(t)
	defer cleanup()

	worker := &domain.Worker{Name: "Vasya", Phone: strPtr("79161112233"), Role: domain.RoleWorker, PasswordHash: "hash"}
	repo.Create(context.Background(), worker)

	body := models.CreateWorkerRequest{
		Name:     "New Worker",
		Phone:    "79009999999",
		Password: "pass123",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/owner/workers", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("UserId", worker.ID.String())

	handler.CreateWorker(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func strPtr(s string) *string { return &s }
