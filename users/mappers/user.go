package mappers

import (
	"image-reports/users/dtos"
	"image-reports/users/models"
)

func MapToUserDTO(user *models.Users) dtos.UserResponse {
	return dtos.UserResponse{
		Id:    user.ID,
		Email: user.Email,
		Role:  user.Role.Name,
	}
}

func MapToUsersDTO(users []models.Users) []dtos.UserResponse {
	res := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		res[i] = MapToUserDTO(&user)
	}
	return res
}
