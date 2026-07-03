package requesthandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// создание заявки с сайта. Публичный.
func (h *RequestHandler) CreateRequest(c *gin.Context) {
	var req models.CreateRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and phone required"})
		return
	}

	request, err := h.requestService.Create(c.Request.Context(), req.Name, req.Phone, req.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// список новых заявок. Только владелец.
func (h *RequestHandler) ListNewRequests(c *gin.Context) {
	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	requests, err := h.requestService.ListNew(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

// все заявки. Только владелец.
func (h *RequestHandler) ListAllRequests(c *gin.Context) {
	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	requests, err := h.requestService.ListAll(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

// создать заказ из заявки. Только владелец.
func (h *RequestHandler) ConvertToOrder(c *gin.Context) {
	var input models.CreateOrderFromRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	order, err := h.requestService.ConvertToOrder(c.Request.Context(), ownerID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// отмена заявки. Только владелец.
func (h *RequestHandler) CancelRequest(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	ownerID, _ := uuid.Parse(c.GetString("UserId"))

	if err := h.requestService.Cancel(c.Request.Context(), ownerID, requestID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "request cancelled"})
}