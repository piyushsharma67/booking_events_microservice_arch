package service

import (
	"context"

	"github.com/piyushsharma67/events_booking/services/auth_service/models"
)

type AuthService interface {
	SignUp(ctx context.Context, user models.CreateUserRequest) (*models.CreateUserRequest, error)
	Login(ctx context.Context, user models.LoginUserRequest) (*models.LoginUserRequest, error)
	Notifier(ctx context.Context, user models.CreateUserRequest) error
}
