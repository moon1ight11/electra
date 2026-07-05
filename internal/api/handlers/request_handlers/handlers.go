package requesthandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// создание заявки с сайта
func (h *RequestHandler) CreateRequest(c *gin.Context) {
	var req models.CreateRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("error in sbjson in create request", "error", err, "req", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and phone required"})
		return
	}

	request, err := h.requestService.Create(c.Request.Context(), req.Name, req.Phone, req.Comment)
	if err != nil {
		h.logger.Error("error in service in create request", "error", err, "req", req)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// список новых заявок
func (h *RequestHandler) ListNewRequests(c *gin.Context) {
	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in list new requests", "error", err, "owner_id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	requests, err := h.requestService.ListNew(c.Request.Context(), ownerID)
	if err != nil {
		h.logger.Error("error in service in list new requests", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

// все заявки
func (h *RequestHandler) ListAllRequests(c *gin.Context) {
	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in list all requests", "error", err, "owner_id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	requests, err := h.requestService.ListAll(c.Request.Context(), ownerID)
	if err != nil {
		h.logger.Error("error in service in list new requests", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

// создать заказ из заявки
func (h *RequestHandler) ConvertToOrder(c *gin.Context) {
	var input models.CreateOrderFromRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("error in sbjson in convert to order", "error", err, "input", input)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in convert to order", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	order, err := h.requestService.ConvertToOrder(c.Request.Context(), ownerID, input)
	if err != nil {
		h.logger.Error("error in service in convert to order", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// отмена заявки
func (h *RequestHandler) CancelRequest(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("error in parce id in cancel request", "error", err, "request_id", requestID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("error in parce id in cancel request", "error", err, "owner_id", ownerID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	if err := h.requestService.Cancel(c.Request.Context(), ownerID, requestID); err != nil {
		h.logger.Error("error in service in cancel request", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "request cancelled"})
}
