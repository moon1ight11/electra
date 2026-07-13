package orderworkerhandlers

import (
	"bytes"
	"context"
	"electra/internal/api/models"
	"electra/internal/config"
	"electra/internal/domain"
	orderworkerservice "electra/internal/services/order_worker_service"
	"electra/internal/storage/repo/database"
	orderworkers "electra/internal/storage/repo/order_workers"
	"electra/internal/storage/repo/orders"
	"electra/internal/storage/repo/workers"
	"electra/pkg/logger"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupOWHandler(t *testing.T) (*OrderWorkerHandler, *workers.WorkerRepo, *orders.OrderRepo, func()) {
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

	orderWorkerRepo := orderworkers.NewOrderWorkerRepo(db)
	workerRepo := workers.NewWorkerRepo(db)
	orderRepo := orders.NewOrderRepo(db)
	svc := orderworkerservice.NewOrderWorkerService(orderWorkerRepo, workerRepo)
	log := logger.NewNopLogger()
	handler := NewOrderWorkerHandler(svc, log)

	return handler, workerRepo, orderRepo, func() {
		db.Close()
	}
}

func TestUpdateReport_Success(t *testing.T) {
	handler, workerRepo, orderRepo, cleanup := setupOWHandler(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	order := &domain.Order{Address: "ул. Ленина", CreatedBy: owner.ID}
	orderRepo.Create(context.Background(), order, []uuid.UUID{owner.ID})

	body := models.UpdateReportInput{
		OrderID:      order.ID,
		TimeSpent:    intPtr(120),
		EarnedAmount: floatPtr(10000),
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/worker/orders/report", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("UserId", owner.ID.String())

	handler.UpdateReport(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "report updated")
}

func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }
func strPtr(s string) *string     { return &s }
