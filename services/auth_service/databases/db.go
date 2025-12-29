package databases

import (
	"context"

	"github.com/piyushsharma67/events_booking/services/auth_service/models"
)

type Database interface {
	InsertUser(ctx context.Context, user *models.UserDocument) error
	GetUserByEmail(ctx context.Context, email string) (*models.UserDocument, error)
}
