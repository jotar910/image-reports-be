package auth

import (
	"fmt"
	"net/http"
	"strings"

	shared_models "image-reports/shared/models"

	"github.com/gin-gonic/gin"
)

type header struct {
	TokenString string `header:"Authorization" binding:"required"`
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := header{}

		if err := c.BindHeader(&h); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "request does not contain an access token"})
			return
		}

		claims, err := ValidateToken(strings.Replace(h.TokenString, "Bearer ", "", 1))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user", claims)
	}
}

func AllowOnlyRole(role shared_models.RolesEnum) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := GetTokenClaim(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if claim.Role != role {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("you must be assigned as %s", role)})
			return
		}
	}
}
