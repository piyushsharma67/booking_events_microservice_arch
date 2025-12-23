package service

import (
	"context"

	"github.com/piyushsharma67/movie_booking/services/auth_service/models"
)

type AuthService interface {
	SignUp(ctx context.Context, user models.User) (models.User, error)
	Login(ctx context.Context, user models.User) (models.User, error)
	Notifier(ctx context.Context,user models.User)error
}
