package repository

// its job is to bring the data from the db
import (
	"context"

	"github.com/piyushsharma67/events_booking/services/auth_service/databases"
	"github.com/piyushsharma67/events_booking/services/auth_service/models"
)

type UserRepository struct {
	db databases.Database
}

func NewUserRepository(db databases.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) InsertUser(ctx context.Context, user *models.UserDocument) error {
	
	return r.db.InsertUser(ctx, user)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserDocument, error) {
	user,err:=r.db.GetUserByEmail(ctx,email)

	return user,err
}
