package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/* type Authorize interface {
	Authorize(ctx context.Context, tokenString string) (*model.Users, error)
} */

type header struct {
	tokenString string `header:"Authorization"`
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := header{}

		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request does not contain an access token"})
			return
		}

		claims, err := ValidateToken(h.tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
