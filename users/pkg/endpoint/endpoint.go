package endpoint

import (
	"errors"
	"net/http"

	"image-reports/users/mappers"
	"image-reports/users/pkg/service"

	shared_dtos "image-reports/shared/dtos"
	user_errors "image-reports/users/errors"

	"github.com/gin-gonic/gin"
)

func CheckCredentials(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials shared_dtos.UserCredentials
		if err := c.BindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		usr, err := svc.CheckCredentials(c, credentials)
		println(usr)
		if err == nil {
			c.JSON(http.StatusOK, mappers.MapToUserDTO(usr))
			return
		}
		if errors.Is(err, user_errors.InvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
