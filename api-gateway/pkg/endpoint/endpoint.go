package endpoint

import (
	"net/http"
	"strconv"

	reporter_dtos "image-reports/api-gateway/dtos/reporter"
	user_dtos "image-reports/api-gateway/dtos/user"
	"image-reports/api-gateway/pkg/service"
	"image-reports/api-gateway/validators"

	"github.com/gin-gonic/gin"
)

func Login(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials user_dtos.UserCredentials
		if err := c.BindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, oerr := svc.Login(c, credentials)
		if oerr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": oerr.Error})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func ListReports(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters reporter_dtos.ListFilters
		if err := c.BindQuery(&filters); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		list, oerr := svc.ListReports(c, filters)
		if oerr != nil {
			c.JSON(oerr.Status, gin.H{"error": oerr.Error})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

func GetReport(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			c.Status(http.StatusNotFound)
			return
		}
		report, oerr := svc.GetReport(c, uint(id))
		if oerr != nil {
			c.JSON(oerr.Status, gin.H{"error": oerr.Error})
			return
		}
		c.JSON(http.StatusOK, report)
	}
}

func CreateReport(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form reporter_dtos.ReportCreation
		if err := c.Bind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validators.ReportCreationValidator(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		report, oerr := svc.CreateReport(c, form)
		if oerr != nil {
			c.JSON(oerr.Status, gin.H{"error": oerr.Error})
			return
		}
		c.JSON(http.StatusOK, report)
	}
}

func ReportApproval(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			c.Status(http.StatusNotFound)
			return
		}
		var reportPatch reporter_dtos.ReportPatch
		if err := c.BindJSON(&reportPatch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		report, oerr := svc.ReportApproval(c, uint(id), reportPatch)
		if oerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": oerr.Error})
			return
		}
		c.JSON(http.StatusOK, report)
	}
}

func GetFile(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, oerr := svc.GetFile(c, c.Param("id"))
		if oerr != nil {
			c.JSON(oerr.Status, gin.H{"error": oerr.Error})
			return
		}
		c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-type"), resp.Body, make(map[string]string))
	}
}
