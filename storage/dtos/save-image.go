package dtos

import (
	"mime/multipart"
)

type SaveImage struct {
	ImageID string                `form:"imageId" binding:"required,max=36"`
	Image   *multipart.FileHeader `form:"image" binding:"required"`
}
