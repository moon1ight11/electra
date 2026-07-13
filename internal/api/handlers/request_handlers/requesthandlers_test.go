package requesthandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"electra/internal/api/models"
	"electra/internal/config"
	"electra/internal/domain"
	"electra/internal/services/request_service"
	"electra/internal/storage/repo/database"
	"electra/internal/storage/repo/orders"
	"electra/internal/storage/repo/requests"
	"electra/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupHandler(t *testing.T) (*RequestHandler, func()) {
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

	requestRepo := requests.NewRequestRepo(db)
	orderRepo := orders.NewOrderRepo(db)
	svc := requestservice.NewRequestService(requestRepo, orderRepo)
	logger := logger.NewNopLogger()
	handler := NewRequestHandler(svc, logger)

	return handler, func() {
		db.Close()
	}
}

func TestCreateRequest_Success(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	body := models.CreateRequestInput{
		Name:    "Иван",
		Phone:   "79991112233",
		Comment: "Нужно подключить дом",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/public/requests", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateRequest(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp domain.Request
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Иван", resp.Name)
	assert.Equal(t, domain.RequestNew, resp.Status)
}

func TestCreateRequest_MissingFields(t *testing.T) {
	handler, cleanup := setupHandler(t)
	defer cleanup()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/public/requests", bytes.NewBuffer([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateRequest(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
