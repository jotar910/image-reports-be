package endpoint

import (
	"net/http"
	"strconv"

	"image-reports/processing/dtos"
	"image-reports/processing/mappers"
	"image-reports/processing/pkg/service"

	"image-reports/helpers/services/auth"
	"image-reports/helpers/services/kafka"
	"image-reports/helpers/validators"

	"github.com/gin-gonic/gin"
)

func AddEvaluation(svc service.Service, message *kafka.ImageProcessedMessage) error {
	_, err := svc.Create(mappers.MapProcessedMessageToEvaluationDTO(message))
	return err
}

func GetEvaluation(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			c.Status(http.StatusNotFound)
			return
		}
		evaluation, err := svc.ReadById(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if evaluation == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "evaluation not found"})
			return
		}
		c.JSON(http.StatusOK, mappers.MapToEvaluationDTO(evaluation))
	}
}

func SearchEvaluations(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query dtos.QuerySearch
		if err := c.BindJSON(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		evaluations, err := svc.ReadAll(query.Ids)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, mappers.MapToEvaluationsDTO(evaluations))
	}
}

func ProcessImage(svc service.Service, maxSize int64, availableExtensions string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := auth.GetTokenClaim(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var form dtos.ProcessImage
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		if err := validators.ImageValidator("image", form.Image, maxSize, availableExtensions); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := svc.Process(claim.Id, form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	}
}
