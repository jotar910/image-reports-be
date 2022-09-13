package service

import (
	"context"

	"image-reports/users/models"

	"image-reports/helpers/services/auth"
	shared_dtos "image-reports/shared/dtos"
	user_errors "image-reports/users/errors"
)

type Service interface {
	CheckCredentials(ctx context.Context, credentials shared_dtos.UserCredentials) (*models.Users, error)
	ReadByEmail(ctx context.Context, email string) (*models.Users, error)
}

type service struct {
}

func NewService() Service {
	return service{}
}

func (svc service) CheckCredentials(ctx context.Context, credentials shared_dtos.UserCredentials) (*models.Users, error) {
	record, err := svc.ReadByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, err
	}

	if record == nil {
		return nil, user_errors.InvalidCredentials
	}
	if encrypted, err := auth.HashPassword(credentials.Password); err != nil {
		return nil, err
	} else if err := auth.PasswordsMatch(encrypted, record.Password); err != nil {
		return nil, user_errors.InvalidCredentials
	}

	return record, nil
}

func (svc service) ReadByEmail(ctx context.Context, email string) (*models.Users, error) {
	return nil, nil
}
