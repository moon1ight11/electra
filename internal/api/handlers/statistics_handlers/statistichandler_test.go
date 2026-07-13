package statisticshandlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/services/statistic_service"
	"electra/internal/storage/repo/database"
	"electra/internal/storage/repo/orders"
	"electra/internal/storage/repo/statistic"
	"electra/internal/storage/repo/workers"
	"electra/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupStatsHandler(t *testing.T) (*StatisticsHandler, *workers.WorkerRepo, *orders.OrderRepo, func()) {
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

	statsRepo := statistic.NewStatisticsRepo(db)
	orderRepo := orders.NewOrderRepo(db)
	workerRepo := workers.NewWorkerRepo(db)
	svc := statisticservice.NewStatisticsService(statsRepo)
	log := logger.NewNopLogger()
	handler := NewStatisticsHandler(svc, log)

	return handler, workerRepo, orderRepo, func() {
		db.Close()
	}
}

func TestSummary_Empty(t *testing.T) {
	handler, _, _, cleanup := setupStatsHandler(t)
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/owner/statistics/summary", nil)

	handler.Summary(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "orders_count")
}

func TestAllWorkers_Empty(t *testing.T) {
	handler, _, _, cleanup := setupStatsHandler(t)
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/owner/statistics/all", nil)
	c.Set("UserId", uuid.New().String())

	handler.AllWorkers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}

func TestByWorker_NotFound(t *testing.T) {
	handler, _, _, cleanup := setupStatsHandler(t)
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(
		http.MethodGet,
		"/api/v1/owner/statistics/workers/"+uuid.New().String()+"?from=2026-01-01&to=2026-12-31",
		nil,
	)
	c.Set("UserId", uuid.New().String())
	c.Params = gin.Params{{Key: "workerId", Value: uuid.New().String()}}

	handler.ByWorker(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "orders_count")
}

func TestByWorker_WithData(t *testing.T) {
	handler, workerRepo, orderRepo, cleanup := setupStatsHandler(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	order := &domain.Order{Address: "ул. Ленина", CreatedBy: owner.ID}
	orderRepo.Create(context.Background(), order, []uuid.UUID{owner.ID})
	orderRepo.Complete(context.Background(), order.ID)

	db, _ := database.PostgresConnection(
		config.Config{
			DataBase: config.DataBaseConfig{
				Host:     "localhost",
				Port:     15432,
				DBName:   "electra",
				User:     "user",
				Password: "user_pass"},
		},
	)
	defer db.Close()
	db.DB.Exec("UPDATE order_workers SET time_spent = 120, earned_amount = 10000 WHERE order_id = $1", order.ID)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(
		http.MethodGet,
		"/api/v1/owner/statistics/workers/"+owner.ID.String()+"?from=2020-01-01&to=2030-12-31",
		nil,
	)
	c.Set("UserId", uuid.New().String())
	c.Params = gin.Params{{Key: "workerId", Value: owner.ID.String()}}

	handler.ByWorker(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "10000")
}

func strPtr(s string) *string { return &s }
