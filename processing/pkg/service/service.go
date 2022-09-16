package service

import (
	"errors"

	"image-reports/processing/dtos"
	"image-reports/processing/mappers"
	"image-reports/processing/models"

	"gorm.io/gorm"
)

type Service interface {
	ReadById(id uint) (*models.Evaluations, error)
	ReadAll(ids []uint) ([]models.Evaluations, error)
	Create(evaluation dtos.Evaluation) (*models.Evaluations, error)
	Process(userId uint, form dtos.ProcessImage) error
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (svc *service) ReadById(id uint) (*models.Evaluations, error) {
	evaluation := new(models.Evaluations)
	tx := svc.db.Preload("Categories").Where(&models.Evaluations{ReportID: id}).Find(evaluation)
	if tx.Error == nil {
		return evaluation, nil
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, tx.Error
}

func (svc *service) ReadAll(ids []uint) ([]models.Evaluations, error) {
	evaluations := make([]models.Evaluations, 0)
	tx := svc.db.Preload("Categories").Where(map[string]interface{}{"report_id": ids}).Find(&evaluations)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return evaluations, nil
}

func (svc *service) Create(evaluation dtos.Evaluation) (*models.Evaluations, error) {
	report := mappers.MapToEvaluations(evaluation)
	if tx := svc.db.Create(report); tx.Error != nil {
		return nil, tx.Error
	}
	return report, nil
}

func (svc *service) Process(userId uint, form dtos.ProcessImage) error {
	go newProcessAlgorithm(form.ReportID, form.ImageID).execute()
	return nil
}
