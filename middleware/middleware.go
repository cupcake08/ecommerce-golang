package middleware

import (
	"net/http"

	"github.com/cupcake08/ecommerce-go/tokens"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "No authorization header",
			})
			c.Abort()
			return
		}
		claims, msg := tokens.ValidateToken(ClientToken)
		if msg == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
