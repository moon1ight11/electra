package statistic

import (
	"context"
	"testing"
	"time"

	"electra/internal/config"
	"electra/internal/storage/repo/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupStats(t *testing.T) (*StatisticsRepo, *database.DataBase, func()) {
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

	repo := NewStatisticsRepo(db)
	return repo, db, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func seedStats(t *testing.T, db *database.DataBase) (workerID, orderID uuid.UUID) {
	t.Helper()

	workerID = uuid.New()
	db.DB.Exec("INSERT INTO workers (id, name, phone, role, password_hash) VALUES ($1, 'Worker', '79001234567', 'worker', 'hash')", workerID)

	ownerID := uuid.New()
	db.DB.Exec("INSERT INTO workers (id, name, phone, role, password_hash) VALUES ($1, 'Owner', '79007654321', 'owner', 'hash')", ownerID)

	orderID = uuid.New()
	db.DB.Exec("INSERT INTO orders (id, address, created_by, completed_at) VALUES ($1, 'ул. Ленина', $2, $3)", orderID, ownerID, time.Now())

	owID := uuid.New()
	db.DB.Exec("INSERT INTO order_workers (id, order_id, worker_id, time_spent, earned_amount) VALUES ($1, $2, $3, 120, 10000)", owID, orderID, workerID)

	return
}

func TestByWorker_NoData(t *testing.T) {
	repo, _, cleanup := setupStats(t)
	defer cleanup()

	row, err := repo.ByWorker(context.Background(), uuid.New(), time.Time{}, time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 0, row.OrdersCount)
}

func TestByWorker_WithData(t *testing.T) {
	repo, db, cleanup := setupStats(t)
	defer cleanup()

	workerID, _ := seedStats(t, db)

	row, err := repo.ByWorker(context.Background(), workerID, time.Time{}, time.Now().Add(24*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, 1, row.OrdersCount)
	assert.Equal(t, 10000.0, row.TotalEarned)
	assert.Equal(t, 120, row.TotalTimeSpent)
}

func TestAllWorkers(t *testing.T) {
	repo, db, cleanup := setupStats(t)
	defer cleanup()

	seedStats(t, db)

	rows, err := repo.AllWorkers(context.Background(), time.Time{}, time.Now().Add(24*time.Hour))
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(rows), 1)
}

func TestSummaryStats(t *testing.T) {
	repo, db, cleanup := setupStats(t)
	defer cleanup()

	seedStats(t, db)

	summary, err := repo.SummaryStats(context.Background(), time.Time{}, time.Now().Add(24*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, 1, summary.OrdersCount)
	assert.Equal(t, 10000.0, summary.TotalEarned)
}
