package models

type EvaluationCategories struct {
	ID             uint   `gorm:"primarykey"`
	EvaluationsID  uint   `gorm:"not null"`
	CategoriesName string `gorm:"not null"`
}
