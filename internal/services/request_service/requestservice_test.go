package requestservice

import (
	"context"
	"electra/internal/api/models"
	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/storage/repo/database"
	"electra/internal/storage/repo/orders"
	"electra/internal/storage/repo/requests"
	"electra/internal/storage/repo/workers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupRequest(t *testing.T) (*RequestService, *orders.OrderRepo, *requests.RequestRepo, *workers.WorkerRepo, func()) {
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

	requestRepo := requests.NewRequestRepo(db)
	orderRepo := orders.NewOrderRepo(db)
	workerRepo := workers.NewWorkerRepo(db)
	service := NewRequestService(requestRepo, orderRepo)

	return service, orderRepo, requestRepo, workerRepo, func() {
		db.DB.Exec("DELETE FROM order_workers")
		db.DB.Exec("DELETE FROM orders")
		db.DB.Exec("DELETE FROM requests")
		db.DB.Exec("DELETE FROM workers")
		db.Close()
	}
}

func TestCreateRequest(t *testing.T) {
	svc, _, _, _, cleanup := setupRequest(t)
	defer cleanup()

	req, err := svc.Create(context.Background(), "Иван", "79991112233", "Подключить дом")
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "Иван", req.Name)
	assert.Equal(t, "79991112233", req.Phone)
	assert.Equal(t, domain.RequestNew, req.Status)
}

func TestCreateRequest_EmptyPhone(t *testing.T) {
	svc, _, _, _, cleanup := setupRequest(t)
	defer cleanup()

	_, err := svc.Create(context.Background(), "Иван", "", "")
	assert.Error(t, err)
}

func TestListNew(t *testing.T) {
	svc, _, _, _, cleanup := setupRequest(t)
	defer cleanup()

	svc.Create(context.Background(), "First", "79161111111", "")
	svc.Create(context.Background(), "Second", "79262222222", "")

	list, err := svc.ListNew(context.Background(), [16]byte{})
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestCancelRequest(t *testing.T) {
	svc, _, reqRepo, _, cleanup := setupRequest(t)
	defer cleanup()

	req, _ := svc.Create(context.Background(), "Иван", "79991112233", "Дом")
	err := svc.Cancel(context.Background(), [16]byte{}, req.ID)
	assert.NoError(t, err)

	req, _ = reqRepo.GetByID(context.Background(), req.ID)
	assert.Equal(t, domain.RequestCancelled, req.Status)
}

func TestCancelRequest_AlreadyProcessed(t *testing.T) {
	svc, _, _, _, cleanup := setupRequest(t)
	defer cleanup()

	req, _ := svc.Create(context.Background(), "Иван", "79991112233", "")
	svc.Cancel(context.Background(), [16]byte{}, req.ID)

	err := svc.Cancel(context.Background(), [16]byte{}, req.ID)
	assert.Error(t, err)
}

func TestConvertToOrder(t *testing.T) {
	svc, _, _, workerRepo, cleanup := setupRequest(t)
	defer cleanup()

	req, _ := svc.Create(context.Background(), "Иван", "79991112233", "Подключить дом")

	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	order, err := svc.ConvertToOrder(context.Background(), owner.ID, models.CreateOrderFromRequestInput{
		RequestID:      req.ID,
		Address:        "ул. Ленина, 5",
		Description:    "Подключение",
		EstimatedPrice: floatPtr(20000),
		PlannedDate:    "2026-07-20",
		WorkerIDs:      []uuid.UUID{owner.ID},
	})
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "ул. Ленина, 5", order.Address)

	reqAfter, _ := svc.requestRepo.GetByID(context.Background(), req.ID)
	assert.Equal(t, domain.RequestConverted, reqAfter.Status)
}

func TestConvertToOrder_AlreadyConverted(t *testing.T) {
	svc, _, _, workerRepo, cleanup := setupRequest(t)
	defer cleanup()

	req, _ := svc.Create(context.Background(), "Иван", "79991112233", "")
	owner := &domain.Worker{Name: "Owner", Phone: strPtr("79001234567"), Role: domain.RoleOwner, PasswordHash: "hash"}
	workerRepo.Create(context.Background(), owner)

	svc.ConvertToOrder(context.Background(), owner.ID, models.CreateOrderFromRequestInput{
		RequestID: req.ID, Address: "ул. Ленина", WorkerIDs: []uuid.UUID{owner.ID},
	})

	_, err := svc.ConvertToOrder(context.Background(), owner.ID, models.CreateOrderFromRequestInput{
		RequestID: req.ID, Address: "ул. Ленина", WorkerIDs: []uuid.UUID{owner.ID},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already processed")
}

func floatPtr(f float64) *float64 { return &f }
func strPtr(s string) *string     { return &s }
