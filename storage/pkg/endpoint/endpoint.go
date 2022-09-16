package endpoint

import (
	"net/http"
	"os"
	"path"

	"image-reports/storage/dtos"

	"image-reports/helpers/validators"

	"github.com/gin-gonic/gin"
)

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

		}
		if err := validators.ImageValidator("image", form.Image, maxSize, availableExtensions); err != nil {
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
