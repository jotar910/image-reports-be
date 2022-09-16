package dtos

import user_dtos "image-reports/api-gateway/dtos/user"

type LoginResponse struct {
	User  user_dtos.UserResponse `json:"user"`
	Token string                 `json:"token"`
}
