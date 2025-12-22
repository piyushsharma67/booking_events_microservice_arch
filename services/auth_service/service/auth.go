package service

import (
	"context"

	"github.com/piyushsharma67/movie_booking/services/auth_service/databases"
	"github.com/piyushsharma67/movie_booking/services/auth_service/models"
	"github.com/piyushsharma67/movie_booking/services/auth_service/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(ctx context.Context, user models.User) (models.User, error)
	Login(ctx context.Context, user models.User) (models.User, error)
}

type authService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) SignUp(ctx context.Context, user models.User) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 10)
	user.PasswordHash = string(hash)
	user.Role = "user"

	userDB:=&databases.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}

	if err := s.repo.InsertUser(ctx, userDB); err != nil {
		return models.User{}, err
	}
	
	return user, nil
}

func (s *authService) Login(ctx context.Context, user models.User) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 10)
	user.PasswordHash = string(hash)
	user.Role = "user"

	userDB, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return models.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash), []byte(user.PasswordHash)); err != nil {
		return models.User{}, err
	}

	userFromDB:=models.User{
		ID:           userDB.ID,
		Name:         userDB.Name,
		Email:        userDB.Email,
		PasswordHash: userDB.PasswordHash,
		Role:         userDB.Role,
		Token:        userDB.Token,
	}

	return userFromDB, nil
}

