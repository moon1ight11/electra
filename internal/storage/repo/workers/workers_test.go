package workers

import (
	"context"
	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/storage/repo/database"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupWorkers(t *testing.T) (*WorkerRepo, func()) {
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

	repo := NewWorkerRepo(db)
	return repo, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func TestCreateWorker(t *testing.T) {
	repo, cleanup := setupWorkers(t)
	defer cleanup()

	w := &domain.Worker{
		Name: "Petya", Phone: strPtr("79007654321"), Role: domain.RoleWorker, PasswordHash: "hash",
	}
	err := repo.Create(context.Background(), w)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, w.ID)
}

func TestGetByID(t *testing.T) {
	repo, cleanup := setupWorkers(t)
	defer cleanup()

	w := &domain.Worker{Name: "Test", Phone: strPtr("79161112233"), Role: domain.RoleWorker, PasswordHash: "hash"}
	repo.Create(context.Background(), w)

	found, err := repo.GetByID(context.Background(), w.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "Test", found.Name)
}

func TestGetByID_NotFound(t *testing.T) {
	repo, cleanup := setupWorkers(t)
	defer cleanup()

	found, err := repo.GetByID(context.Background(), uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestGetByPhone(t *testing.T) {
	repo, cleanup := setupWorkers(t)
	defer cleanup()

	w := &domain.Worker{Name: "Test", Phone: strPtr("79161112233"), Role: domain.RoleWorker, PasswordHash: "hash"}
	repo.Create(context.Background(), w)

	found, err := repo.GetByPhone(context.Background(), "79161112233")
	assert.NoError(t, err)
	assert.NotNil(t, found)
}

func TestGetByPhone_NotFound(t *testing.T) {
	repo, cleanup := setupWorkers(t)
	defer cleanup()

	found, err := repo.GetByPhone(context.Background(), "79999999999")
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestList(t *testing.T) {
	repo, cleanup := setupWorkers(t)
	defer cleanup()

	repo.Create(context.Background(), &domain.Worker{Name: "A", Phone: strPtr("79161111111"), Role: domain.RoleWorker, PasswordHash: "h"})
	repo.Create(context.Background(), &domain.Worker{Name: "B", Phone: strPtr("79262222222"), Role: domain.RoleWorker, PasswordHash: "h"})

	list, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func strPtr(s string) *string { return &s }
