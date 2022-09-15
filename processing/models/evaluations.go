package models

import "gorm.io/gorm"

type Evaluations struct {
	gorm.Model
	ReportID   uint         `gorm:"not null"`
	ImageID    string       `gorm:"type:varchar(36);not null"`
	Grade      int          `gorm:"not null;min=0;max=0"`
	Categories []Categories `gorm:"many2many:evaluation_categories"`
}
