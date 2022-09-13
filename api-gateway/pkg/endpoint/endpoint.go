package endpoint

import (
	"net/http"

	"image-reports/api-gateway/pkg/service"

	shared_dtos "image-reports/shared/dtos"

	"github.com/gin-gonic/gin"
)

func Login(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials shared_dtos.UserCredentials
		if err := c.BindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		resp, err := svc.Login(c, credentials)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
