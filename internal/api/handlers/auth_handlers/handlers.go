package authhandlers

import (
	"electra/internal/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Логин по телефону и паролю
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and password required"})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("cookie", token, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged in"})
}

// Логаут
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("cookie", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
