package middleware

import (
	"net/http"
	"tiketsepur/utils"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "user role tidak ditemukan", nil)
			c.Abort()
			return
		}

		role := userRole.(string)

		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ErrorResponse(c, http.StatusForbidden, "kamu tidak memiliki izin untuk mengakses ini", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

func UserOrAdmin() gin.HandlerFunc {
	return RoleMiddleware("user", "admin")
}