package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/piyushsharma67/movie_booking/services/auth_service/databases"
	"github.com/piyushsharma67/movie_booking/services/auth_service/models"
	"github.com/piyushsharma67/movie_booking/services/auth_service/repository"
	"github.com/piyushsharma67/movie_booking/services/auth_service/utils"
)



type authService struct {
	repo *repository.UserRepository
	notifier Notifier
}

func NewAuthService(repo *repository.UserRepository,notifier Notifier) AuthService {
	return &authService{repo: repo,notifier: notifier}
}

func (s *authService) SignUp(ctx context.Context, user models.User) (models.User, error) {
	// 1. Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return models.User{}, err
	}

	// 2. Prepare DB model
	userDB := &databases.User{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	// 3. Insert into DB
	if err := s.repo.InsertUser(ctx, userDB); err != nil {
		return models.User{}, err
	}

	// 4. Generate JWT
	token, err := utils.GenerateJWT(
		userDB.ID,
		userDB.Email,
		userDB.Role,
		os.Getenv("JWT_SECRET"),
	)
	if err != nil {
		return models.User{}, err
	}

	// 5. Return API response (NO password/hash)
	return models.User{
		ID:    userDB.ID,
		Name:  userDB.Name,
		Email: userDB.Email,
		Role:  userDB.Role,
		Token: token,
	}, nil
}

func (s *authService) Login(ctx context.Context, user models.User) (models.User, error) {
	// 1. Fetch user from DB
	userDB, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return models.User{}, err
	}

	// 2. Compare password
	if err := utils.CheckPassword(user.Password, userDB.PasswordHash); err != nil {
		return models.User{}, errors.New("invalid credentials")
	}

	// 3. Generate JWT
	token, err := utils.GenerateJWT(
		userDB.ID,
		userDB.Email,
		userDB.Role,
		os.Getenv("JWT_SECRET"),
	)
	if err != nil {
		return models.User{}, err
	}

	// 4. Return response (NO password/hash)
	return models.User{
		ID:    userDB.ID,
		Name:  userDB.Name,
		Email: userDB.Email,
		Role:  userDB.Role,
		Token: token,
	}, nil
}

func (s *authService)Notifier(ctx context.Context,user models.User)error{
	return s.notifier.SendNotification(user.Email,"Welcome",fmt.Sprintf("Hi %s Welome to the Booking Application", user.Name))
}
