package dtos

import (
	shared_models "image-reports/shared/models"

	"mime/multipart"
)

type ReportCreation struct {
	Name     string                           `form:"name" binding:"required,max=255"`
	Callback string                           `form:"callback" binding:"required,max=2048"`
	Type     shared_models.ReportCreationType `form:"type" binding:"required"`
	Url      string                           `form:"url" binding:"max=2048"`
	File     *multipart.FileHeader            `form:"file"`
}

type ReportCreationData struct {
	Name     string `json:"name"`
	Callback string `json:"callback"`
	ImageID  string `json:"imageId"`
}

type SaveImage struct {
	ImageID string                `form:"imageId"`
	Image   *multipart.FileHeader `form:"image"`
}

type ProcessImage struct {
	ReportID uint                  `form:"reportId"`
	ImageID  string                `form:"imageId"`
	Image    *multipart.FileHeader `form:"image"`
}
