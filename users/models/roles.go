package models

import (
	shared_models "image-reports/shared/models"
)

type Roles struct {
	ID   uint                    `gorm:"primaryKey"`
	Name shared_models.RolesEnum `gorm:"type:varchar(50);unique;not null"`
}
