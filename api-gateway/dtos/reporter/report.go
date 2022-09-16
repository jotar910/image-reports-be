package dtos

import (
	processing_dtos "image-reports/api-gateway/dtos/processing"

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
	Date     int64                          `json:"date"`
}

type ReportOutbound struct {
	ID         uint                                `json:"id"`
	Name       string                              `json:"name"`
	User       string                              `json:"user"`
	Image      string                              `json:"image"`
	Status     shared_models.ReportStatusEnum      `json:"status"`
	Approval   *Approval                           `json:"approval"`
	Evaluation *processing_dtos.EvaluationOutbound `json:"evaluation"`
	Date       int64                               `json:"creationDate"`
}
