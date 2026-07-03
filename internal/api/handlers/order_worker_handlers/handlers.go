package orderworkerhandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdateReport — заполнить/обновить отчёт по заказу.
func (h *OrderWorkerHandler) UpdateReport(c *gin.Context) {
	var input models.UpdateReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	workerID, _ := uuid.Parse(c.GetString("UserId"))

	if err := h.orderWorkerService.UpdateReport(c.Request.Context(), workerID, input); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "report updated"})
}

// GetByOrder — отчёты всех исполнителей по заказу.
func (h *OrderWorkerHandler) GetByOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("UserId"))

	reports, err := h.orderWorkerService.GetByOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// RemoveWorker — снять исполнителя с заказа. Только владелец.
func (h *OrderWorkerHandler) RemoveWorker(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	workerID, err := uuid.Parse(c.Param("workerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	if err := h.orderWorkerService.RemoveWorker(c.Request.Context(), ownerID, orderID, workerID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "worker removed"})
}

// ListCompleted — история выполненных заказов исполнителя.
func (h *OrderWorkerHandler) ListCompleted(c *gin.Context) {
	workerID, _ := uuid.Parse(c.GetString("UserId"))

	orders, err := h.orderWorkerService.ListCompletedByWorker(c.Request.Context(), workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// ListAllCompleted — все выполненные заказы. Только владелец.
func (h *OrderWorkerHandler) ListAllCompleted(c *gin.Context) {
	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	orders, err := h.orderWorkerService.ListAllCompleted(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}