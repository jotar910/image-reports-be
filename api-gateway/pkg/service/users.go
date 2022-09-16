package service

import (
	"context"
	"image-reports/api-gateway/dtos"
	user_dtos "image-reports/api-gateway/dtos/user"
)

func (svc service) GetUserById(ctx context.Context, id uint) (*user_dtos.UserResponse, *dtos.ErrorOutbound) {
	return svc.userClient.GetUserById(ctx, id)
}
