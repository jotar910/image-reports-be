package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Email    string `gorm:"type:varchar(255);not null"`
	Password string `gorm:"type:varchar(500);not null"`
	RoleId   uint   `gorm:"not null"`
	Role     Roles  `gorm:"not null"`
}
