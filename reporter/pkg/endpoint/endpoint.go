package endpoint

import (
	"context"
	"net/http"
	"strconv"

	"image-reports/reporter/dtos"
	"image-reports/reporter/mappers"
	"image-reports/reporter/models"
	"image-reports/reporter/pkg/service"

	"image-reports/helpers/services/auth"
	"image-reports/helpers/services/kafka"

	shared_models "image-reports/shared/models"

	"github.com/gin-gonic/gin"
)

func OnImageProcessedMessage(
	ctx context.Context,
	message *kafka.ImageProcessedMessage,
	svc service.Service,
) error {
	var err error
	if message.Err == nil {
		if message.Going {
			_, err = svc.PatchStatus(message.ReportId, shared_models.ReportStatusEvaluating)

		} else {
			_, err = svc.PatchGrade(message.ReportId, message.Grade)
		}
	} else {
		_, err = svc.PatchStatus(message.ReportId, shared_models.ReportStatusError)
	}
	return err
}

func ListReports(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters := dtos.ListFilters{}
		if err := c.BindQuery(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filters.CheckDefaults()
		reports, err := svc.ReadAll(filters)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		total, err := svc.Count()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, mappers.MapToReportsListDTO(reports, total, filters.Page, filters.Count))
	}
}

func GetReport(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			c.Status(http.StatusNotFound)
			return
		}
		var report *models.Reports
		report, err = svc.ReadById(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if report == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
			return
		}
		c.JSON(http.StatusOK, mappers.MapToReportDTO(report))
	}
}

func CreateReport(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := auth.GetTokenClaim(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var reportCreation dtos.ReportCreation
		if err := c.BindJSON(&reportCreation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		report, err := svc.Create(claim.Id, reportCreation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, mappers.MapToReportDTO(report))
	}
}

func ReportApproval(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := auth.GetTokenClaim(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			c.Status(http.StatusNotFound)
			return
		}
		var reportPatch dtos.ReportPatch
		if err := c.BindJSON(&reportPatch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		report, err := svc.PatchApproval(uint(id), claim.Id, reportPatch.ApprovalStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if report == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
			return
		}
		c.JSON(http.StatusOK, mappers.MapToReportDTO(report))
	}
}
