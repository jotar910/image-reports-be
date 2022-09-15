package dtos

import (
	"mime/multipart"
)

type ProcessImage struct {
	ReportID uint                  `form:"reportId" binding:"required"`
	ImageID  string                `form:"imageId" binding:"required,max=36"`
	Image    *multipart.FileHeader `form:"image" binding:"required"`
}
