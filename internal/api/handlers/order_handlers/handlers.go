package orderhandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *OrderHandler) CreateDirect(c *gin.Context) {
	var input models.CreateOrderDirectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("failed to bind create order request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if len(input.WorkerIDs) == 0 {
		h.logger.Error("empty worker_ids in create order")
		c.JSON(http.StatusBadRequest, gin.H{"error": "выберите хотя бы одного исполнителя"})
		return
	}

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	order, err := h.orderService.CreateDirect(c.Request.Context(), ownerID, input)
	if err != nil {
		h.logger.Error("failed to create order", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("order created", "order_id", order.ID)

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) ListPlanned(c *gin.Context) {
	workerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse worker id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	orders, err := h.orderService.ListPlannedByWorker(c.Request.Context(), workerID)
	if err != nil {
		h.logger.Error("failed to list planned orders", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) ListAllPlanned(c *gin.Context) {
	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	orders, err := h.orderService.ListAllPlanned(c.Request.Context(), ownerID)
	if err != nil {
		h.logger.Error("failed to list all planned orders", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) Complete(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("failed to parse order id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	workerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse worker id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
		return
	}

	if err := h.orderService.Complete(c.Request.Context(), workerID, orderID); err != nil {
		h.logger.Error("failed to complete order", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("order completed", "order_id", orderID)

	c.JSON(http.StatusOK, gin.H{"message": "order completed"})
}

func (h *OrderHandler) CompleteByOwner(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("failed to parse order id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	if err := h.orderService.CompleteByOwner(c.Request.Context(), orderID); err != nil {
		h.logger.Error("failed to complete order by owner", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("order completed by owner", "order_id", orderID)

	c.JSON(http.StatusOK, gin.H{"message": "order completed"})
}

func (h *OrderHandler) Update(c *gin.Context) {
	var input models.UpdateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("failed to bind update order request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	order, err := h.orderService.Update(c.Request.Context(), input)
	if err != nil {
		h.logger.Error("failed to update order", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("order updated", "order_id", input.ID)

	c.JSON(http.StatusOK, order)
}