package requesthandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *RequestHandler) CreateRequest(c *gin.Context) {
	var req models.CreateRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind create request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and phone required"})
		return
	}

	request, err := h.requestService.Create(c.Request.Context(), req.Name, req.Phone, req.Comment)
	if err != nil {
		h.logger.Error("failed to create request", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("request created from landing", "req_id", request.ID)

	c.JSON(http.StatusCreated, request)
}

func (h *RequestHandler) ConvertToOrder(c *gin.Context) {
	var input models.CreateOrderFromRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("failed to bind convert to order request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	order, err := h.requestService.ConvertToOrder(c.Request.Context(), ownerID, input)
	if err != nil {
		h.logger.Error("failed to convert request to order", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("order created from request", "order_id", order.ID, "req_id", input.RequestID)

	c.JSON(http.StatusCreated, order)
}

func (h *RequestHandler) ListNewRequests(c *gin.Context) {
	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	requests, err := h.requestService.ListNew(c.Request.Context(), ownerID)
	if err != nil {
		h.logger.Error("failed to list new requests", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (h *RequestHandler) ListAllRequests(c *gin.Context) {
	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	requests, err := h.requestService.ListAll(c.Request.Context(), ownerID)
	if err != nil {
		h.logger.Error("failed to list all requests", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (h *RequestHandler) CancelRequest(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("failed to parse request id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	ownerID, err := uuid.Parse(c.GetString("UserId"))
	if err != nil {
		h.logger.Error("failed to parse owner id", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid owner Id"})
		return
	}

	if err := h.requestService.Cancel(c.Request.Context(), ownerID, requestID); err != nil {
		h.logger.Error("failed to cancel request", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("request cancelled", "req_id", requestID)

	c.JSON(http.StatusOK, gin.H{"message": "request cancelled"})
}