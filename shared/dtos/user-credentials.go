package dtos

type UserCredentials struct {
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,max=100"`
}
