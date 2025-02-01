package middleware

import (
	"be-test/helpers"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}
		fmt.Println(token)

		claims, err := helpers.ValidateJWT(strings.Split(token, " ")[1])
		if err != nil {
			if err.Error() == "Token is expired" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			}
			c.Abort()
			return
		}

		c.Set("user", claims["email"])
		c.Next()
	}
}
