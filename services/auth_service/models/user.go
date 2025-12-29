package models

type CreateUserRequest struct {
	ID string `json:"_id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	PasswordHash string `json:"password_hash"`
	Role     string `json:"role"`
	Token string `json:"token"`
}

type LoginUserRequest struct {
	ID string `json:"_id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	PasswordHash string `json:"password_hash"`
	Token string `json:"token"`
	Role string `json:"role"`
}
