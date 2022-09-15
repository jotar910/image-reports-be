package models

type Categories struct {
	Name string `gorm:"primaryKey;type:varchar(255);not null"`
}
