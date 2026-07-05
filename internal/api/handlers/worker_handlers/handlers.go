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

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in create worker", "error", err, "owner_id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	worker, err := h.authService.CreateWorker(c.Request.Context(), ownerID, req.Name, req.Phone, req.Password)
	if err != nil {
		h.logger.Error("error in service in create worker", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, worker)
}
