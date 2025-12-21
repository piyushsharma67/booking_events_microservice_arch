package databases

import (
	"context"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         string
}

type Database interface {
	InsertUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
}
