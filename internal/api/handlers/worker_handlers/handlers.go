package workerhandlers

import (
	"electra/internal/api/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *WorkerHandler) CreateWorker(c *gin.Context) {
	var req models.CreateWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind create worker request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, phone and password required"})
		return
	}

	ownerID := c.GetString("UserId")
	if ownerID == "" {
		h.logger.Error("empty owner id in create worker")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	worker, err := h.authService.CreateWorker(c.Request.Context(), ownerUUID, req.Name, req.Phone, req.Password, req.Specialization)
	if err != nil {
		h.logger.Error("failed to create worker", "error", err)
		if strings.Contains(err.Error(), "only owner") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("worker created", "worker_id", worker.ID)

	c.JSON(http.StatusCreated, worker)
}

func (h *WorkerHandler) ListWorkers(c *gin.Context) {
	workers, err := h.authService.ListWorkers(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list workers", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workers)
}

func (h *WorkerHandler) GetMe(c *gin.Context) {
	userIDStr := c.GetString("UserId")
	if userIDStr == "" {
		h.logger.Error("empty user id in get me")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Error("failed to parse user id", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	worker, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get me", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, worker)
}

func (h *WorkerHandler) DeleteWorker(c *gin.Context) {
	workerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("failed to parse worker id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.authService.DeleteWorker(c.Request.Context(), ownerID, workerID); err != nil {
		h.logger.Error("failed to delete worker", "error", err)
		if strings.Contains(err.Error(), "only owner") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("worker deleted", "worker_id", workerID)

	c.JSON(http.StatusOK, gin.H{"message": "worker deleted"})
}

func (h *WorkerHandler) UpdateProfile(c *gin.Context) {
	var input models.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("failed to bind update profile request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	userID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse user id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	worker, err := h.authService.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
		h.logger.Error("failed to update profile", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("profile updated", "user_id", userID)

	c.JSON(http.StatusOK, worker)
}