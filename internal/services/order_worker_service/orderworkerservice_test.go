package orderworkerservice

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

func setupOW(t *testing.T) (*OrderWorkerService, *workers.WorkerRepo, *orders.OrderRepo, func()) {
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

	orderWorkerRepo := orderworkers.NewOrderWorkerRepo(db)
	workerRepo := workers.NewWorkerRepo(db)
	orderRepo := orders.NewOrderRepo(db)
	service := NewOrderWorkerService(orderWorkerRepo, workerRepo)

	return service, workerRepo, orderRepo, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func createOrderWithWorkers(t *testing.T, orderRepo *orders.OrderRepo, workerRepo *workers.WorkerRepo) (*domain.Order, *domain.Worker, *domain.Worker) {
	t.Helper()
	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	worker := &domain.Worker{Name: "Petya", Phone: strPtr("79007654321"), Role: domain.RoleWorker, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), worker)

	order := &domain.Order{
		Address:   "ул. Ленина, 5",
		CreatedBy: owner.ID,
	}
	orderRepo.Create(context.Background(), order, []uuid.UUID{owner.ID, worker.ID})

	return order, owner, worker
}

func TestUpdateReport(t *testing.T) {
	svc, _, orderRepo, cleanup := setupOW(t)
	defer cleanup()

	order, owner, _ := createOrderWithWorkers(t, orderRepo, svc.workerRepo)

	err := svc.UpdateReport(context.Background(), owner.ID, models.UpdateReportInput{
		OrderID:       order.ID,
		TimeSpent:     intPtr(120),
		EarnedAmount:  floatPtr(10000),
		MaterialsUsed: strPtr("Кабель 30м"),
		Notes:         strPtr("Готово"),
	})
	assert.NoError(t, err)

	reports, _ := svc.GetByOrder(context.Background(), owner.ID, order.ID)
	assert.Len(t, reports, 2)
}

func TestUpdateReport_NotAssigned(t *testing.T) {
	svc, _, orderRepo, cleanup := setupOW(t)
	defer cleanup()

	order, _, _ := createOrderWithWorkers(t, orderRepo, svc.workerRepo)

	stranger := &domain.Worker{Name: "Stranger", Phone: strPtr("79163334455"), Role: domain.RoleWorker, PasswordHash: "hash"}
	svc.workerRepo.Create(context.Background(), stranger)

	err := svc.UpdateReport(context.Background(), stranger.ID, models.UpdateReportInput{
		OrderID: order.ID,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not assigned")
}

func TestGetByOrder_OwnerAccess(t *testing.T) {
	svc, _, orderRepo, cleanup := setupOW(t)
	defer cleanup()

	order, owner, _ := createOrderWithWorkers(t, orderRepo, svc.workerRepo)

	reports, err := svc.GetByOrder(context.Background(), owner.ID, order.ID)
	assert.NoError(t, err)
	assert.Len(t, reports, 2)
}

func TestRemoveWorker(t *testing.T) {
	svc, _, orderRepo, cleanup := setupOW(t)
	defer cleanup()

	order, _, worker := createOrderWithWorkers(t, orderRepo, svc.workerRepo)

	err := svc.RemoveWorker(context.Background(), uuid.Nil, order.ID, worker.ID)
	assert.NoError(t, err)

	reports, _ := svc.GetByOrder(context.Background(), worker.ID, order.ID)
	assert.Len(t, reports, 0)
}

func TestListCompletedByWorker(t *testing.T) {
	svc, _, orderRepo, cleanup := setupOW(t)
	defer cleanup()

	order, owner, _ := createOrderWithWorkers(t, orderRepo, svc.workerRepo)
	orderRepo.Complete(context.Background(), order.ID)

	orders, err := svc.ListCompletedByWorker(context.Background(), owner.ID)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
}

func TestListAllCompleted(t *testing.T) {
	svc, _, orderRepo, cleanup := setupOW(t)
	defer cleanup()

	order, _, _ := createOrderWithWorkers(t, orderRepo, svc.workerRepo)
	orderRepo.Complete(context.Background(), order.ID)

	orders, err := svc.ListAllCompleted(context.Background(), uuid.Nil)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
}

func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }
func strPtr(s string) *string     { return &s }
