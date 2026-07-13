package workerhandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// создание нового исполнителя
func (h *WorkerHandler) CreateWorker(c *gin.Context) {
	var req models.CreateWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("error in sbjson in create worker", "error", err, "req", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, phone and password required"})
		return
	}

	ownerID := c.GetString("UserId")
	if ownerID == "" {
		h.logger.Error("error in get owner id in create worker")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		h.logger.Error("error in parce id in create worker", "error", err, "owner_id", ownerID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	worker, err := h.authService.CreateWorker(c.Request.Context(), ownerUUID, req.Name, req.Phone, req.Password, req.Specialization)
	if err != nil {
		h.logger.Error("error in services in create worker", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("worker created successfully", "worker_id", worker.ID)

	c.JSON(http.StatusCreated, worker)
}

// получение списка работников
func (h *WorkerHandler) ListWorkers(c *gin.Context) {
	workers, err := h.authService.ListWorkers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workers)
}

// получение информации о текущем пользователе
func (h *WorkerHandler) GetMe(c *gin.Context) {
	userIDStr := c.GetString("UserId")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	worker, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, worker)
}
