package middleware

import (
	"net/http"
	"strconv"

	Token "crm-app/backend/token"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		err := Token.TokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		data, err := Token.ExtractTokenID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		// c.AddParam("userId", strconv.FormatUint(uint64(data), 10))
		user_id := strconv.FormatUint(uint64(data), 10)
		c.Set("userId", user_id)
		c.Next()
	}
}
