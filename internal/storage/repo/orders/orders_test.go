package orders

import (
	"context"
	"testing"

	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/storage/repo/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupOrders(t *testing.T) (*OrderRepo, *database.DataBase, func()) {
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

	db.DB.Exec("INSERT INTO workers (id, name, phone, role, password_hash) VALUES ($1, 'Owner', '79001234567', 'owner', 'hash')", uuid.New())

	repo := NewOrderRepo(db)
	return repo, db, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func getOwnerID(t *testing.T, db *database.DataBase) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	err := db.DB.QueryRow("SELECT id FROM workers WHERE role = 'owner' LIMIT 1").Scan(&id)
	if err != nil {
		t.Fatalf("no owner: %v", err)
	}
	return id
}

func getWorkerID(t *testing.T, db *database.DataBase) uuid.UUID {
	t.Helper()
	id := uuid.New()
	db.DB.Exec("INSERT INTO workers (id, name, phone, role, password_hash) VALUES ($1, 'Worker', '79007654321', 'worker', 'hash')", id)
	return id
}

func TestCreateOrder(t *testing.T) {
	repo, db, cleanup := setupOrders(t)
	defer cleanup()

	ownerID := getOwnerID(t, db)
	workerID := getWorkerID(t, db)

	order := &domain.Order{
		Address:   "ул. Ленина, 5",
		CreatedBy: ownerID,
	}
	err := repo.Create(context.Background(), order, []uuid.UUID{ownerID, workerID})
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, order.ID)
}

func TestGetByID(t *testing.T) {
	repo, db, cleanup := setupOrders(t)
	defer cleanup()

	ownerID := getOwnerID(t, db)
	order := &domain.Order{Address: "ул. Ленина, 5", CreatedBy: ownerID}
	repo.Create(context.Background(), order, []uuid.UUID{ownerID})

	found, err := repo.GetByID(context.Background(), order.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "ул. Ленина, 5", found.Address)
}

func TestGetByID_NotFound(t *testing.T) {
	repo, _, cleanup := setupOrders(t)
	defer cleanup()

	found, err := repo.GetByID(context.Background(), uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestListPlannedByWorker(t *testing.T) {
	repo, db, cleanup := setupOrders(t)
	defer cleanup()

	ownerID := getOwnerID(t, db)
	order := &domain.Order{Address: "ул. Ленина", CreatedBy: ownerID}
	repo.Create(context.Background(), order, []uuid.UUID{ownerID})

	list, err := repo.ListPlannedByWorker(context.Background(), ownerID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestListAllPlanned(t *testing.T) {
	repo, db, cleanup := setupOrders(t)
	defer cleanup()

	ownerID := getOwnerID(t, db)
	repo.Create(context.Background(), &domain.Order{Address: "A", CreatedBy: ownerID}, []uuid.UUID{ownerID})
	repo.Create(context.Background(), &domain.Order{Address: "B", CreatedBy: ownerID}, []uuid.UUID{ownerID})

	list, err := repo.ListAllPlanned(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestComplete(t *testing.T) {
	repo, db, cleanup := setupOrders(t)
	defer cleanup()

	ownerID := getOwnerID(t, db)
	order := &domain.Order{Address: "ул. Ленина", CreatedBy: ownerID}
	repo.Create(context.Background(), order, []uuid.UUID{ownerID})

	err := repo.Complete(context.Background(), order.ID)
	assert.NoError(t, err)

	found, _ := repo.GetByID(context.Background(), order.ID)
	assert.NotNil(t, found.CompletedAt)
}

func TestComplete_AlreadyCompleted(t *testing.T) {
	repo, db, cleanup := setupOrders(t)
	defer cleanup()

	ownerID := getOwnerID(t, db)
	order := &domain.Order{Address: "ул. Ленина", CreatedBy: ownerID}
	repo.Create(context.Background(), order, []uuid.UUID{ownerID})
	repo.Complete(context.Background(), order.ID)

	err := repo.Complete(context.Background(), order.ID)
	assert.Error(t, err)
}

func TestUpdate(t *testing.T) {
	repo, db, cleanup := setupOrders(t)
	defer cleanup()

	ownerID := getOwnerID(t, db)
	order := &domain.Order{Address: "ул. Старая", CreatedBy: ownerID}
	repo.Create(context.Background(), order, []uuid.UUID{ownerID})

	order.Address = "ул. Новая"
	err := repo.Update(context.Background(), order)
	assert.NoError(t, err)

	found, _ := repo.GetByID(context.Background(), order.ID)
	assert.Equal(t, "ул. Новая", found.Address)
}
