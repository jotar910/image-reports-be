package endpoint

import (
	"context"
	"net/http"
	"os"
	"path"
	"time"

	"image-reports/storage/dtos"

	"image-reports/helpers/services/kafka"
	"image-reports/helpers/validators"

	"github.com/gin-gonic/gin"
)

func OnReportCreatedMessage(ctx context.Context, message *kafka.ReportCreatedMessage) (*kafka.ImageStoredMessage, error) {
	time.Sleep(time.Second * 200)
	return kafka.NewImageStoredMessageCompleted(message.ReportId, message.ImageId), nil
}

func OnReportDeletedMessage(ctx context.Context, message *kafka.ReportDeletedMessage) error {
	time.Sleep(time.Second * 200)
	return nil
}

func GetImage(folder string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File(path.Join(folder, c.Param("id")))
	}
}

func SaveImage(folder string, maxSize int64, availableExtensions string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form dtos.SaveImage
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		} else if err := validators.ImageValidator("image", form.Image, maxSize, availableExtensions); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := c.SaveUploadedFile(form.Image, path.Join(folder, form.ImageID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	}
}
