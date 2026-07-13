package statisticservice

import (
	"context"
	"electra/internal/config"
	"electra/internal/storage/repo/database"
	"electra/internal/storage/repo/orders"
	"electra/internal/storage/repo/statistic"
	"electra/internal/storage/repo/workers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupStats(t *testing.T) (*StatisticsService, *orders.OrderRepo, *workers.WorkerRepo, func()) {
	t.Helper()

	db, err := database.PostgresConnection(config.Config{
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

	statsRepo := statistic.NewStatisticsRepo(db)
	orderRepo := orders.NewOrderRepo(db)
	workerRepo := workers.NewWorkerRepo(db)
	service := NewStatisticsService(statsRepo)

	return service, orderRepo, workerRepo, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func TestByWorker_NoData(t *testing.T) {
	svc, _, _, cleanup := setupStats(t)
	defer cleanup()

	stats, err := svc.ByWorker(context.Background(), uuid.New(), uuid.New(), "2026-01-01", "2026-12-31")
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 0, stats.OrdersCount)
}

func TestAllWorkers_Empty(t *testing.T) {
	svc, _, _, cleanup := setupStats(t)
	defer cleanup()

	stats, err := svc.AllWorkers(context.Background(), uuid.New(), "2026-01-01", "2026-12-31")
	assert.NoError(t, err)
	assert.Empty(t, stats)
}

func TestSummary_NoData(t *testing.T) {
	svc, _, _, cleanup := setupStats(t)
	defer cleanup()

	summary, err := svc.Summary(context.Background(), "2026-01-01", "2026-12-31")
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 0, summary.OrdersCount)
}