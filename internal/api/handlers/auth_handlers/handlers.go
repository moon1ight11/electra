package authhandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind login request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and password required"})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		h.logger.Error("failed to login", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("cookie", token, 3600, "/", "", true, true)

	h.logger.Info("user logged in", "user_phone", req.Phone)

	c.JSON(http.StatusOK, gin.H{"message": "logged in"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("cookie", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
