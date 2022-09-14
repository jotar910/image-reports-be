package models

import (
	shared_models "image-reports/shared/models"

	"gorm.io/gorm"
)

type Approvals struct {
	gorm.Model
	ReportID uint                             `gorm:"not null"`
	UserID   uint                             `gorm:"not null"`
	Status   shared_models.ApprovalStatusEnum `gorm:"type:varchar(50);unique;not null"`
}
