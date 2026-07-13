package orderhandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"electra/internal/api/models"
	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/services/order_service"
	"electra/internal/storage/repo/database"
	orderworkers "electra/internal/storage/repo/order_workers"
	"electra/internal/storage/repo/orders"
	"electra/internal/storage/repo/workers"
	"electra/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupOrderHandler(t *testing.T) (*OrderHandler, *workers.WorkerRepo, *orders.OrderRepo, func()) {
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

	orderRepo := orders.NewOrderRepo(db)
	orderWorkerRepo := orderworkers.NewOrderWorkerRepo(db)
	workerRepo := workers.NewWorkerRepo(db)
	orderSvc := orderservice.NewOrderService(orderRepo, orderWorkerRepo)
	log := logger.NewNopLogger()
	handler := NewOrderHandler(orderSvc, log)

	return handler, workerRepo, orderRepo, func() {
		db.Close()
	}
}

func createOwner(t *testing.T, repo *workers.WorkerRepo) *domain.Worker {
	t.Helper()
	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	repo.Create(context.Background(), owner)
	return owner
}

func TestCreateDirect_Success(t *testing.T) {
	handler, workerRepo, _, cleanup := setupOrderHandler(t)
	defer cleanup()

	owner := createOwner(t, workerRepo)

	body := models.CreateOrderDirectInput{
		Address:        "ул. Ленина, 5",
		Description:    "Подключение",
		EstimatedPrice: floatPtr(20000),
		PlannedDate:    "2026-07-20",
		WorkerIDs:      []uuid.UUID{owner.ID},
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/owner/orders/direct", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("UserId", owner.ID.String())
	c.Set("UserRole", domain.RoleOwner)

	handler.CreateDirect(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateDirect_EmptyWorkers(t *testing.T) {
	handler, workerRepo, _, cleanup := setupOrderHandler(t)
	defer cleanup()

	owner := createOwner(t, workerRepo)

	body := models.CreateOrderDirectInput{
		Address:   "ул. Ленина, 5",
		WorkerIDs: []uuid.UUID{},
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/owner/orders/direct", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("UserId", owner.ID.String())
	c.Set("UserRole", domain.RoleOwner)

	handler.CreateDirect(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func floatPtr(f float64) *float64 { return &f }

func strPtr(s string) *string { return &s }
