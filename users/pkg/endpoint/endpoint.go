package endpoint

import (
	"errors"
	"net/http"
	"strconv"

	"image-reports/users/dtos"
	"image-reports/users/mappers"
	"image-reports/users/pkg/service"

	user_errors "image-reports/users/errors"

	"github.com/gin-gonic/gin"
)

func CheckCredentials(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials dtos.UserCredentials
		if err := c.BindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		usr, err := svc.CheckCredentials(c, credentials)
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

func CheckUserId(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			c.Status(http.StatusNotFound)
			return
		}
		usr, err := svc.ReadById(c, uint(id))
		if err != nil || usr == nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusOK)
	}
}

func GetUser(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			c.Status(http.StatusNotFound)
			return
		}
		user, err := svc.ReadById(c, uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, mappers.MapToUserDTO(user))
	}
}

func SearchUsers(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query dtos.QuerySearch
		if err := c.BindJSON(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		users, err := svc.ReadAll(c, query.Ids)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, mappers.MapToUsersDTO(users))
	}
}
