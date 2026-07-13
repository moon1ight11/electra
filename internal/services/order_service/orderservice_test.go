package orderservice

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/storage/repo/database"
	orderworkers "electra/internal/storage/repo/order_workers"
	"electra/internal/storage/repo/orders"
	"electra/internal/storage/repo/workers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupOrder(t *testing.T) (*OrderService, *orders.OrderRepo, *workers.WorkerRepo, func()) {
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

	orderRepo := orders.NewOrderRepo(db)
	orderWorkerRepo := orderworkers.NewOrderWorkerRepo(db)
	workerRepo := workers.NewWorkerRepo(db)
	service := NewOrderService(orderRepo, orderWorkerRepo)

	return service, orderRepo, workerRepo, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func TestCreateDirect(t *testing.T) {
	svc, _, workerRepo, cleanup := setupOrder(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	order, err := svc.CreateDirect(context.Background(), owner.ID, models.CreateOrderDirectInput{
		Address:        "ул. Пушкина, 10",
		Description:    "Замена щитка",
		EstimatedPrice: floatPtr(15000),
		PlannedDate:    "2026-07-25",
		WorkerIDs:      []uuid.UUID{owner.ID},
	})
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "ул. Пушкина, 10", order.Address)
}

func TestListPlannedByWorker(t *testing.T) {
	svc, _, workerRepo, cleanup := setupOrder(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	svc.CreateDirect(context.Background(), owner.ID, models.CreateOrderDirectInput{
		Address: "ул. Ленина, 5", WorkerIDs: []uuid.UUID{owner.ID},
	})

	orders, err := svc.ListPlannedByWorker(context.Background(), owner.ID)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
}

func TestListAllPlanned(t *testing.T) {
	svc, _, workerRepo, cleanup := setupOrder(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	svc.CreateDirect(context.Background(), owner.ID, models.CreateOrderDirectInput{
		Address: "ул. Ленина, 5", WorkerIDs: []uuid.UUID{owner.ID},
	})
	svc.CreateDirect(context.Background(), owner.ID, models.CreateOrderDirectInput{
		Address: "ул. Пушкина, 10", WorkerIDs: []uuid.UUID{owner.ID},
	})

	orders, err := svc.ListAllPlanned(context.Background(), owner.ID)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestComplete(t *testing.T) {
	svc, _, workerRepo, cleanup := setupOrder(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	order, _ := svc.CreateDirect(context.Background(), owner.ID, models.CreateOrderDirectInput{
		Address: "ул. Ленина, 5", WorkerIDs: []uuid.UUID{owner.ID},
	})

	err := svc.Complete(context.Background(), owner.ID, order.ID)
	assert.NoError(t, err)
}

func TestComplete_NotAssigned(t *testing.T) {
	svc, _, workerRepo, cleanup := setupOrder(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	worker := &domain.Worker{Name: "Worker", Phone: strPtr("79007654321"), Role: domain.RoleWorker, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), worker)

	order, _ := svc.CreateDirect(context.Background(), owner.ID, models.CreateOrderDirectInput{
		Address: "ул. Ленина, 5", WorkerIDs: []uuid.UUID{owner.ID},
	})

	err := svc.Complete(context.Background(), worker.ID, order.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not assigned")
}

func TestComplete_AlreadyCompleted(t *testing.T) {
	svc, _, workerRepo, cleanup := setupOrder(t)
	defer cleanup()

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	order, _ := svc.CreateDirect(context.Background(), owner.ID, models.CreateOrderDirectInput{
		Address: "ул. Ленина, 5", WorkerIDs: []uuid.UUID{owner.ID},
	})

	svc.Complete(context.Background(), owner.ID, order.ID)

	err := svc.Complete(context.Background(), owner.ID, order.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already completed")
}

func floatPtr(f float64) *float64 { return &f }
func strPtr(s string) *string     { return &s }
