package models

type CreateUserRequest struct {
	ID string `json:"name"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	PasswordHash string `json:"password_hash"`
	Role     string `json:"role" binding:"required"`
	Token string `json:"token"`
}
