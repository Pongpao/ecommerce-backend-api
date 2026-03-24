package middleware

import (
	"strings"
	
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"project-e-commerce/config"
	"project-e-commerce/utils"
)

var jwtSecret = []byte(config.GetEnv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Error(utils.Unauthorized("missing token", nil))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.Error(utils.Unauthorized("invalid token", nil))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Error(utils.Unauthorized("invalid token claims", nil))
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		c.Next()
	}
}