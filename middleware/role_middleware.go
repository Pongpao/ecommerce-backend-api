package middleware

import (
	"project-e-commerce/utils"

	"github.com/gin-gonic/gin"
)

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleValue, exists := c.Get("role")
		if !exists {
			c.Error(utils.Forbidden("role not found", nil))
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.Error(utils.Forbidden("invalid role", nil))
			c.Abort()
			return
		}

		if role != requiredRole {
			c.Error(utils.Forbidden("access denied", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}