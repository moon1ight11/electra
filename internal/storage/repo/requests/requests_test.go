package requests

import (
	"context"
	"testing"

	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/storage/repo/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupRequests(t *testing.T) (*RequestRepo, func()) {
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

	repo := NewRequestRepo(db)
	return repo, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func TestCreateRequest(t *testing.T) {
	repo, cleanup := setupRequests(t)
	defer cleanup()

	req := &domain.Request{Name: "Иван", Phone: "79991112233"}
	err := repo.Create(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, req.ID)
	assert.Equal(t, domain.RequestNew, req.Status)
}

func TestGetByID(t *testing.T) {
	repo, cleanup := setupRequests(t)
	defer cleanup()

	req := &domain.Request{Name: "Иван", Phone: "79991112233"}
	repo.Create(context.Background(), req)

	found, err := repo.GetByID(context.Background(), req.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "Иван", found.Name)
}

func TestGetByID_NotFound(t *testing.T) {
	repo, cleanup := setupRequests(t)
	defer cleanup()

	found, err := repo.GetByID(context.Background(), uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestList(t *testing.T) {
	repo, cleanup := setupRequests(t)
	defer cleanup()

	repo.Create(context.Background(), &domain.Request{Name: "A", Phone: "79161111111"})
	repo.Create(context.Background(), &domain.Request{Name: "B", Phone: "79262222222"})

	list, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestListNew(t *testing.T) {
	repo, cleanup := setupRequests(t)
	defer cleanup()

	repo.Create(context.Background(), &domain.Request{Name: "New", Phone: "79161111111"})
	req := &domain.Request{Name: "Old", Phone: "79262222222"}
	repo.Create(context.Background(), req)
	repo.MarkConverted(context.Background(), req.ID)

	list, err := repo.ListNew(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "New", list[0].Name)
}

func TestMarkConverted(t *testing.T) {
	repo, cleanup := setupRequests(t)
	defer cleanup()

	req := &domain.Request{Name: "Иван", Phone: "79991112233"}
	repo.Create(context.Background(), req)

	err := repo.MarkConverted(context.Background(), req.ID)
	assert.NoError(t, err)

	found, _ := repo.GetByID(context.Background(), req.ID)
	assert.Equal(t, domain.RequestConverted, found.Status)
}

func TestMarkCancelled(t *testing.T) {
	repo, cleanup := setupRequests(t)
	defer cleanup()

	req := &domain.Request{Name: "Иван", Phone: "79991112233"}
	repo.Create(context.Background(), req)

	err := repo.MarkCancelled(context.Background(), req.ID)
	assert.NoError(t, err)

	found, _ := repo.GetByID(context.Background(), req.ID)
	assert.Equal(t, domain.RequestCancelled, found.Status)
}
