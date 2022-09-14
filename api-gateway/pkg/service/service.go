package service

import (
	"context"

	"image-reports/api-gateway/dtos"

	"image-reports/helpers/services/auth"
	users_client "image-reports/helpers/services/http-client/users"
	shared_dtos "image-reports/shared/dtos"
)

type Service interface {
	Login(ctx context.Context, credentials shared_dtos.UserCredentials) (*dtos.LoginResponse, error)
}

type service struct {
	usersClient users_client.HttpClient
}

func NewService(usersClient users_client.HttpClient) Service {
	return &service{usersClient}
}

func (svc service) Login(ctx context.Context, credentials shared_dtos.UserCredentials) (*dtos.LoginResponse, error) {
	usr, err := svc.usersClient.CheckCredentials(credentials)
	if err != nil {
		return nil, err
	}

	tokenString, err := auth.GenerateJWT(usr.Id, usr.Email, usr.Role)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginResponse{
		User:  usr,
		Token: tokenString,
	}, nil
}
