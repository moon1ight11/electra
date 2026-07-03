package orderhandlers

import (
	"electra/internal/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// CreateDirect — создать заказ без заявки. Только владелец.
func (h *OrderHandler) CreateDirect(c *gin.Context) {
	var input models.CreateOrderDirectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	order, err := h.orderService.CreateDirect(c.Request.Context(), ownerID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// ListPlanned — запланированные заказы текущего исполнителя.
func (h *OrderHandler) ListPlanned(c *gin.Context) {
	workerID, _ := uuid.Parse(c.GetString("UserId"))

	orders, err := h.orderService.ListPlannedByWorker(c.Request.Context(), workerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// ListAllPlanned — все запланированные заказы. Только владелец.
func (h *OrderHandler) ListAllPlanned(c *gin.Context) {
	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	orders, err := h.orderService.ListAllPlanned(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// Complete — завершить заказ. Исполнитель или владелец.
func (h *OrderHandler) Complete(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	workerID, _ := uuid.Parse(c.GetString("UserId"))

	if err := h.orderService.Complete(c.Request.Context(), workerID, orderID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order completed"})
}
