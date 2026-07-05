package orderhandlers

import (
	"electra/internal/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// создание заказа мимо заявки
func (h *OrderHandler) CreateDirect(c *gin.Context) {
	var input models.CreateOrderDirectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("error in sbjson in create direct", "error", err, "input", input)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in create direct", "error", err, "id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	order, err := h.orderService.CreateDirect(c.Request.Context(), ownerID, input)
	if err != nil {
		h.logger.Error("error in service in create direct", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// запланированные заказы текущего исполнителя
func (h *OrderHandler) ListPlanned(c *gin.Context) {
	workerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in list planned", "error", err, "id", workerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	orders, err := h.orderService.ListPlannedByWorker(c.Request.Context(), workerID)
	if err != nil {
		h.logger.Error("error in service in list planned", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// все запланированные заказы
func (h *OrderHandler) ListAllPlanned(c *gin.Context) {
	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in list all planned", "error", err, "id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	orders, err := h.orderService.ListAllPlanned(c.Request.Context(), ownerID)
	if err != nil {
		h.logger.Error("error in service in list all planned", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// завершить заказ
func (h *OrderHandler) Complete(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("error in parce id in order complete", "error", err, "order_id", orderID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	workerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in order complete", "error", err, "worker_id", workerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	if err := h.orderService.Complete(c.Request.Context(), workerID, orderID); err != nil {
		h.logger.Error("error in service in order complete", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order completed"})
}
