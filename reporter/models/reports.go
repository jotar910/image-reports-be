package models

import (
	shared_models "image-reports/shared/models"

	"gorm.io/gorm"
)

type Reports struct {
	gorm.Model
	Name     string                         `gorm:"type:varchar(255);not null"`
	UserID   uint                           `gorm:"not null"`
	ImageID  string                         `gorm:"type:varchar(36);not null"`
	Status   shared_models.ReportStatusEnum `gorm:"not null;default:'NEW'"`
	Approval Approvals                      `gorm:"foreignKey:ReportID"`
}
