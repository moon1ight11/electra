package orderworkers

import (
	"context"
	"testing"

	"electra/internal/config"
	"electra/internal/storage/repo/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupOW(t *testing.T) (*OrderWorkerRepo, *database.DataBase, func()) {
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

	repo := NewOrderWorkerRepo(db)
	return repo, db, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func seedWorkers(t *testing.T, db *database.DataBase) (ownerID, workerID uuid.UUID) {
	t.Helper()
	ownerID = uuid.New()
	workerID = uuid.New()
	db.DB.Exec("INSERT INTO workers (id, name, phone, role, password_hash) VALUES ($1, 'Owner', '79001234567', 'owner', 'hash')", ownerID)
	db.DB.Exec("INSERT INTO workers (id, name, phone, role, password_hash) VALUES ($1, 'Worker', '79007654321', 'worker', 'hash')", workerID)
	return
}

func seedOrder(t *testing.T, db *database.DataBase, ownerID uuid.UUID, workerIDs []uuid.UUID) uuid.UUID {
	t.Helper()
	orderID := uuid.New()
	db.DB.Exec("INSERT INTO orders (id, address, created_by) VALUES ($1, 'ул. Ленина', $2)", orderID, ownerID)
	for _, wid := range workerIDs {
		db.DB.Exec("INSERT INTO order_workers (id, order_id, worker_id) VALUES ($1, $2, $3)", uuid.New(), orderID, wid)
	}
	return orderID
}

func TestGetByOrder(t *testing.T) {
	repo, db, cleanup := setupOW(t)
	defer cleanup()

	ownerID, workerID := seedWorkers(t, db)
	orderID := seedOrder(t, db, ownerID, []uuid.UUID{ownerID, workerID})

	list, err := repo.GetByOrder(context.Background(), orderID)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetByOrderAndWorker(t *testing.T) {
	repo, db, cleanup := setupOW(t)
	defer cleanup()

	ownerID, workerID := seedWorkers(t, db)
	orderID := seedOrder(t, db, ownerID, []uuid.UUID{ownerID, workerID})

	ow, err := repo.GetByOrderAndWorker(context.Background(), orderID, workerID)
	assert.NoError(t, err)
	assert.NotNil(t, ow)
	assert.Equal(t, workerID, ow.WorkerID)
}

func TestGetByOrderAndWorker_NotFound(t *testing.T) {
	repo, db, cleanup := setupOW(t)
	defer cleanup()

	ownerID, _ := seedWorkers(t, db)
	orderID := seedOrder(t, db, ownerID, []uuid.UUID{ownerID})

	ow, err := repo.GetByOrderAndWorker(context.Background(), orderID, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, ow)
}

func TestUpdate(t *testing.T) {
	repo, db, cleanup := setupOW(t)
	defer cleanup()

	ownerID, _ := seedWorkers(t, db)
	orderID := seedOrder(t, db, ownerID, []uuid.UUID{ownerID})

	list, _ := repo.GetByOrder(context.Background(), orderID)
	ow := list[0]
	timeSpent := 120
	earned := 10000.0
	materials := "Кабель"
	notes := "Готово"

	ow.TimeSpent = &timeSpent
	ow.EarnedAmount = &earned
	ow.MaterialsUsed = &materials
	ow.Notes = &notes

	err := repo.Update(context.Background(), &ow)
	assert.NoError(t, err)

	updated, _ := repo.GetByOrderAndWorker(context.Background(), orderID, ow.WorkerID)
	assert.Equal(t, 120, *updated.TimeSpent)
	assert.Equal(t, 10000.0, *updated.EarnedAmount)
}

func TestRemove(t *testing.T) {
	repo, db, cleanup := setupOW(t)
	defer cleanup()

	ownerID, workerID := seedWorkers(t, db)
	orderID := seedOrder(t, db, ownerID, []uuid.UUID{ownerID, workerID})

	err := repo.Remove(context.Background(), orderID, workerID)
	assert.NoError(t, err)

	list, _ := repo.GetByOrder(context.Background(), orderID)
	assert.Len(t, list, 1)
}

func TestListCompletedByWorker(t *testing.T) {
	repo, db, cleanup := setupOW(t)
	defer cleanup()

	ownerID, _ := seedWorkers(t, db)
	orderID := seedOrder(t, db, ownerID, []uuid.UUID{ownerID})
	db.DB.Exec("UPDATE orders SET completed_at = NOW() WHERE id = $1", orderID)

	list, err := repo.ListCompletedByWorker(context.Background(), ownerID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestListAllCompleted(t *testing.T) {
	repo, db, cleanup := setupOW(t)
	defer cleanup()

	ownerID, _ := seedWorkers(t, db)
	orderID := seedOrder(t, db, ownerID, []uuid.UUID{ownerID})
	db.DB.Exec("UPDATE orders SET completed_at = NOW() WHERE id = $1", orderID)

	list, err := repo.ListAllCompleted(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}
