package middleware

import (
	"net/http"
	"project-e-commerce/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		if len(c.Errors) > 0 {

			err := c.Errors.Last().Err

			if appErr, ok := err.(*utils.AppError); ok {
				c.JSON(appErr.Code, gin.H{
					"success": false,
					"message": appErr.Message,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "internal server error",
			})
		}
	}
}