package orderworkerhandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// заполнить/обновить отчёт по заказу
func (h *OrderWorkerHandler) UpdateReport(c *gin.Context) {
	var input models.UpdateReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("error in sbjson in update report", "error", err, "input", input)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	workerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in update report", "error", err, "id", workerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	if err := h.orderWorkerService.UpdateReport(c.Request.Context(), workerID, input); err != nil {
		h.logger.Error("error in service in update report", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "report updated"})
}

// отчёты всех исполнителей по заказу
func (h *OrderWorkerHandler) GetByOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("error in parce id in get by order", "error", err, "order_id", orderID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	userID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in get by order", "error", err, "user_id", userID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	reports, err := h.orderWorkerService.GetByOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		h.logger.Error("error in service in get by order", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// снять исполнителя с заказа
func (h *OrderWorkerHandler) RemoveWorker(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		h.logger.Error("error in parce id in remowe worker", "error", err, "order_id", orderID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	workerID, err := uuid.Parse(c.Param("workerId"))
	if err != nil {
		h.logger.Error("error in parce id in remowe worker", "error", err, "worker_id", workerID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in remowe worker", "error", err, "owner_id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	if err := h.orderWorkerService.RemoveWorker(c.Request.Context(), ownerID, orderID, workerID); err != nil {
		h.logger.Error("error in service in remowe worker", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "worker removed"})
}

// история выполненных заказов исполнителя
func (h *OrderWorkerHandler) ListCompleted(c *gin.Context) {
	workerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in list completed", "error", err, "worker_id", workerID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid worker id"})
		return
	}

	orders, err := h.orderWorkerService.ListCompletedByWorker(c.Request.Context(), workerID)
	if err != nil {
		h.logger.Error("error in service in list completed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// все выполненные заказы
func (h *OrderWorkerHandler) ListAllCompleted(c *gin.Context) {
	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in list all completed", "error", err, "owner_id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	orders, err := h.orderWorkerService.ListAllCompleted(c.Request.Context(), ownerID)
	if err != nil {
		h.logger.Error("error in service in list all completed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
