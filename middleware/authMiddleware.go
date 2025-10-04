package middleware

import (
	"net/http"
	"strings"
	"tiketsepur/service"
	"tiketsepur/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}
	
	claims, err := authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "tidak ada token atau expired token", err)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("token", tokenString)

		c.Next()
	}
}