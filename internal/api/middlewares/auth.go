package middlewares

import (
	"electra/internal/api/jwt"
	"electra/internal/domain"
	"electra/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// аус-мидлвар работника
func Auth(jwtService jwt.TokenService, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, err := c.Cookie("cookie")
		if err != nil {
			logger.Error("Error in Auth-middleware", "error:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		claims := jwt.Claims{}
		token, err := jwtService.ParseToken(value, &claims)
		if err != nil {
			logger.Error("Error in Auth-middleware", "error:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		if !token.Valid {
			logger.Error("Error in Auth-middleware", "error:", "token not valid")
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Set("UserId", *claims.UserId)
		c.Set("UserRole", claims.Role)

		c.Next()
	}
}

// аус-мидлвар для владельца
func AuthOwner(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("UserRole")
		if !exists {
			logger.Error("Error in AuthOwner", "error:", "role not found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		r, ok := role.(domain.WorkerRole)
		if !ok || r != domain.RoleOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "owner only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
