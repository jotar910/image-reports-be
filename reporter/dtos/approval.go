package dtos

import (
	shared_models "image-reports/shared/models"
)

type Approval struct {
	UserID uint                             `json:"userId"`
	Status shared_models.ApprovalStatusEnum `json:"status"`
	Date   string                           `json:"date"`
}
