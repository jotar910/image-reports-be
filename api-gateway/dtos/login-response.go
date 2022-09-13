package dtos

import shared_dtos "image-reports/shared/dtos"

type LoginResponse struct {
	User  shared_dtos.UserResponse `json:"user"`
	Token string                   `json:"token"`
}
