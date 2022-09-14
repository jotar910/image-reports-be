package service

import (
	"context"
	"errors"

	"image-reports/users/models"

	"image-reports/helpers/services/auth"
	shared_dtos "image-reports/shared/dtos"
	user_errors "image-reports/users/errors"

	"gorm.io/gorm"
)

type Service interface {
	CheckCredentials(ctx context.Context, credentials shared_dtos.UserCredentials) (*models.Users, error)
	ReadByEmail(ctx context.Context, email string) (*models.Users, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return service{db}
}

func (svc service) CheckCredentials(ctx context.Context, credentials shared_dtos.UserCredentials) (*models.Users, error) {
	record, err := svc.ReadByEmail(ctx, credentials.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.InvalidCredentials
		}
		return nil, err
	}

	if err := auth.PasswordsMatch(record.Password, credentials.Password); err != nil {
		return nil, user_errors.InvalidCredentials
	}

	return record, nil
}

func (svc service) ReadByEmail(ctx context.Context, email string) (*models.Users, error) {
	user := new(models.Users)
	tx := svc.db.Joins("Role").Where(&models.Users{Email: email}).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
