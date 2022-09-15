package service

import (
	"errors"
	"fmt"

	"image-reports/reporter/dtos"
	"image-reports/reporter/models"

	shared_models "image-reports/shared/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service interface {
	Count() (int64, error)
	ReadAll(dtos.ListFilters) ([]*models.Reports, error)
	ReadById(uint) (*models.Reports, error)
	Create(uint, dtos.ReportCreation) (*models.Reports, error)
	PatchStatus(uint, shared_models.ReportStatusEnum) (*models.Reports, error)
	PatchGrade(uint, int) (*models.Reports, error)
	PatchApproval(id, userId uint, status shared_models.ApprovalStatusEnum) (*models.Reports, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

func (svc *service) Count() (int64, error) {
	var count int64
	if tx := svc.db.Model(&models.Reports{}).Select("id").Count(&count); tx.Error != nil {
		return 0, tx.Error
	}
	return count, nil
}

func (svc *service) ReadAll(filters dtos.ListFilters) ([]*models.Reports, error) {
	reports := make([]*models.Reports, 0)
	tx := svc.db.
		Joins("Approval").
		Offset(int((filters.Page - 1) * filters.Count)).
		Limit(int(filters.Count)).
		Order("grade ASC NULLS LAST, created_at DESC").
		Find(&reports)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return reports, nil
}

func (svc *service) ReadById(id uint) (*models.Reports, error) {
	report := new(models.Reports)
	tx := svc.db.Joins("Approval").First(report, id)
	if tx.Error == nil {
		return report, nil
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, tx.Error
}

func (svc *service) Create(userId uint, creation dtos.ReportCreation) (*models.Reports, error) {
	report := &models.Reports{
		Name:     creation.Name,
		UserID:   userId,
		ImageID:  creation.ImageID,
		Callback: creation.Callback,
	}
	if tx := svc.db.Create(report); tx.Error != nil {
		return nil, tx.Error
	}
	return report, nil
}

func (svc *service) PatchStatus(id uint, status shared_models.ReportStatusEnum) (*models.Reports, error) {
	report := &models.Reports{Model: gorm.Model{ID: id}}
	tx := svc.db.Model(report).Clauses(clause.Returning{}).Update("status", status)
	if tx.Error == nil {
		return report, nil
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, tx.Error
}

func (svc *service) PatchGrade(id uint, grade int) (*models.Reports, error) {
	report := &models.Reports{Model: gorm.Model{ID: id}}
	tx := svc.db.Model(report).Clauses(clause.Returning{}).Update("grade", grade)
	if tx.Error == nil {
		return report, nil
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, tx.Error
}

func (svc *service) PatchApproval(id, userId uint, status shared_models.ApprovalStatusEnum) (*models.Reports, error) {
	report := new(models.Reports)
	tx := svc.db.Begin()
	tx = tx.Joins("Approval").First(report, id)
	if tx.Error != nil {
		tx.Rollback()
		fmt.Printf("%v\n", tx.Error)
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	if report.Status != shared_models.ReportStatusPending {
		return nil, errors.New("report must be in a pending state")
	}
	report.Status = shared_models.ReportStatusPublished
	report.Approval = models.Approvals{ReportID: id, UserID: userId, Status: status}
	tx = tx.Updates(report)
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}
	tx.Commit()
	return report, nil
}
