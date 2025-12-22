package repository

import (
	"context"

	"github.com/piyushsharma67/movie_booking/services/auth_service/databases"
	"github.com/piyushsharma67/movie_booking/services/auth_service/models"
)

type UserRepository struct {
	db databases.Database
}

func NewUserRepository(db databases.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) InsertUser(ctx context.Context, user *databases.User) error {
	userDB:=&models.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}
	return r.db.InsertUser(ctx, userDB)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*databases.User, error) {
	user,err:=r.db.GetUserByEmail(ctx,email)

	if err!=nil{
		return nil,err
	}
	userDB := &databases.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}
	return userDB,nil
}
