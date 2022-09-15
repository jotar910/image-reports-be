package dtos

import (
	shared_models "image-reports/shared/models"
)

type ReportPatch struct {
	ApprovalStatus shared_models.ApprovalStatusEnum `json:"approvalStatus"`
}
