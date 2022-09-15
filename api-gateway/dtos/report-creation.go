package dtos

import (
	shared_models "image-reports/shared/models"

	"mime/multipart"
)

type ReportCreation struct {
	Name     string                           `form:"name" binding:"required,max=255"`
	Callback string                           `form:"callback" binding:"required,max=2048"`
	Type     shared_models.ReportCreationType `form:"type" binding:"required"`
	Url      string                           `form:"url" binding:"urlTypeRequired,max=2048"`
	File     *multipart.FileHeader            `form:"file" binding:"fileTypeRequired,max=2048"`
}
