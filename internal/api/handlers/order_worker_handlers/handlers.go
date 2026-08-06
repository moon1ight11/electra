package orderworkerhandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *OrderWorkerHandler) UpdateReport(c *gin.Context) {
	var input models.UpdateReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("failed to bind update report request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	workerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse worker id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	if err := h.orderWorkerService.UpdateReport(c.Request.Context(), workerID, input); err != nil {
		h.logger.Error("failed to update report", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("report updated", "order_id", input.OrderID)

	c.JSON(http.StatusOK, gin.H{"message": "report updated"})
}

func (h *OrderWorkerHandler) GetByOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("failed to parse order id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	userID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse user id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	reports, err := h.orderWorkerService.GetByOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		h.logger.Error("failed to get reports", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (h *OrderWorkerHandler) ListCompleted(c *gin.Context) {
	workerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse worker id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	orders, err := h.orderWorkerService.ListCompletedByWorker(c.Request.Context(), workerID)
	if err != nil {
		h.logger.Error("failed to list completed orders", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderWorkerHandler) ListAllCompleted(c *gin.Context) {
	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	orders, err := h.orderWorkerService.ListAllCompleted(c.Request.Context(), ownerID)
	if err != nil {
		h.logger.Error("failed to list all completed orders", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderWorkerHandler) RemoveWorker(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		h.logger.Error("failed to parse order id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	workerID, err := uuid.Parse(c.Param("workerId"))
	if err != nil {
		h.logger.Error("failed to parse worker id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	if err := h.orderWorkerService.RemoveWorker(c.Request.Context(), ownerID, orderID, workerID); err != nil {
		h.logger.Error("failed to remove worker", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("worker removed from order", "order_id", orderID, "worker_id", workerID)

	c.JSON(http.StatusOK, gin.H{"message": "worker removed"})
}