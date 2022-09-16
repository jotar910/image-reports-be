package transport

import (
	"image-reports/api-gateway/pkg/service"
	"image-reports/helpers/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUserValidity(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := auth.GetTokenClaim(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if ok := svc.CheckUserById(c, claim.Id); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
			return
		}
	}
}

func AddContextToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("token", c.Request.Header.Get("Authorization"))
	}
}
