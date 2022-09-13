package mappers

import (
	shared_dtos "image-reports/shared/dtos"
	"image-reports/users/models"
)

func MapToUserDTO(user *models.Users) shared_dtos.UserResponse {
	return shared_dtos.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}
}
