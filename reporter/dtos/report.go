package dtos

import (
	shared_models "image-reports/shared/models"
)

type Report struct {
	ID       uint                           `json:"id"`
	Name     string                         `json:"name"`
	UserID   uint                           `json:"userId"`
	ImageID  string                         `json:"imageId"`
	Callback string                         `json:"callback"`
	Status   shared_models.ReportStatusEnum `json:"status"`
	Approval *Approval                      `json:"approval"`
	Date     string                         `json:"date"`
}
