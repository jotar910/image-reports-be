package service

import (
	"context"

	"image-reports/api-gateway/dtos"
	user_dtos "image-reports/api-gateway/dtos/user"

	"image-reports/helpers/services/auth"
)

func (svc service) Login(ctx context.Context, credentials user_dtos.UserCredentials) (*dtos.LoginResponse, *dtos.ErrorOutbound) {
	usr, oerr := svc.userClient.CheckCredentials(credentials)
	if oerr != nil {
		return nil, oerr
	}

	tokenString, err := auth.GenerateJWT(usr.Id, usr.Email, usr.Role)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}

	return &dtos.LoginResponse{
		User:  *usr,
		Token: tokenString,
	}, nil
}

func (svc service) CheckUserById(ctx context.Context, id uint) bool {
	return svc.userClient.CheckUserById(ctx, id) == nil
}
