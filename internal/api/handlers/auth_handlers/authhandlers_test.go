package authhandlers

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

func setupAuthHandler(t *testing.T) (*AuthHandler, *workers.WorkerRepo, func()) {
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
	handler := NewAuthHandler(authService, log)

	return handler, workerRepo, func() {
		db.Close()
	}
}

func TestLogin_Success(t *testing.T) {
	handler, repo, cleanup := setupAuthHandler(t)
	defer cleanup()

	hash, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.DefaultCost)
	repo.Create(context.Background(), &domain.Worker{
		Name: "Test", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: string(hash),
	})

	body := models.LoginRequest{Phone: "79001234567", Password: "mypassword"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "logged in")
}

func TestLogin_InvalidPassword(t *testing.T) {
	handler, repo, cleanup := setupAuthHandler(t)
	defer cleanup()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	repo.Create(context.Background(), &domain.Worker{
		Name: "Test", Phone: strPtr("79001234567"), Role: domain.RoleWorker, PasswordHash: string(hash),
	})

	body := models.LoginRequest{Phone: "79001234567", Password: "wrong"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func strPtr(s string) *string { return &s }
