package endpoint

import (
	"context"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"image-reports/storage/dtos"

	"image-reports/helpers/services/auth"
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
		claim, err := auth.GetTokenClaim(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.File(path.Join(folder, strconv.Itoa(int(claim.Id)), c.Param("id")))
	}
}

func SaveImage(folder string, maxSize int64, availableExtensions string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := auth.GetTokenClaim(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var form dtos.SaveImage
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		if err := validators.ImageValidator("image", form.Image, maxSize, availableExtensions); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		savingFolder := path.Join(folder, strconv.Itoa(int(claim.Id)))
		if err := os.MkdirAll(savingFolder, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := c.SaveUploadedFile(form.Image, path.Join(savingFolder, form.ImageID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	}
}
