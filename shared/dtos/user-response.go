package dtos

import (
	shared_models "image-reports/shared/models"
)

type UserResponse struct {
	Id    uint                    `json:"id"`
	Email string                  `json:"email"`
	Role  shared_models.RolesEnum `json:"role"`
}
