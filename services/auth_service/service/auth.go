package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/piyushsharma67/events_booking/services/auth_service/logger"
	"github.com/piyushsharma67/events_booking/services/auth_service/models"
	"github.com/piyushsharma67/events_booking/services/auth_service/repository"
	"github.com/piyushsharma67/events_booking/services/auth_service/utils"
)

type authService struct {
	repo     *repository.UserRepository
	notifier Notifier
	logger   logger.Logger
}

func NewAuthService(repo *repository.UserRepository, notifier Notifier, logger logger.Logger) AuthService {
	return &authService{repo: repo, notifier: notifier, logger: logger}
}

func (s *authService) SignUp(ctx context.Context, user models.CreateUserRequest) (*models.CreateUserRequest, error) {

	// 1. Hash password
	log := s.logger.WithContext(ctx)

	log.Info("Signup started",
		"user", user.Name,
	)

	//checking if the user with same email already exists

	userdb,err:=s.repo.GetUserByEmail(ctx,user.Email)

	if err!=nil{
		return nil,err
	}else if userdb.Email!=""{
		return nil,errors.New(utils.USER_ALREADY_EXISTS)
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = hashedPassword
	userDB, err := models.MapCreateRequestToDocument(&user)

	if err != nil {
		return nil, err
	}

	// 3. Insert into DB
	if err := s.repo.InsertUser(ctx, userDB); err != nil {
		return nil, err
	}

	// 4. Generate JWT
	token, err := utils.GenerateJWT(
		userDB.ID.String(),
		userDB.Email,
		userDB.Role,
		os.Getenv("JWT_SECRET"),
	)
	if err != nil {
		return nil, err
	}

	// 5. Return API response (NO password/hash)
	return &models.CreateUserRequest{
		ID:    userDB.ID.Hex(),
		Name:  userDB.Name,
		Email: userDB.Email,
		Role:  userDB.Role,
		Token: token,
	}, nil
}

func (s *authService) Login(ctx context.Context, user models.LoginUserRequest) (*models.LoginUserRequest, error) {
	reqID, _ := ctx.Value("request_id").(string)

	s.logger.Info("login request",
		"request_id", reqID,
		"user", user.Email,
	)
	// 1. Fetch user from DB
	userDB, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	// 2. Compare password
	if err := utils.CheckPassword(user.Password, userDB.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3. Generate JWT
	token, err := utils.GenerateJWT(
		userDB.ID.String(),
		userDB.Email,
		userDB.Role,
		os.Getenv("JWT_SECRET"),
	)
	if err != nil {
		return nil, err
	}

	err = s.notifier.SendNotification(userDB.Email, "Login Alert", fmt.Sprintf("Hi %s, You have successfully logged in.", userDB.Name))

	if err != nil {
		return nil, nil
	}
	// 4. Return response (NO password/hash)
	return &models.LoginUserRequest{
		ID:    userDB.ID.Hex(),
		Email: userDB.Email,
		Token: token,
		Role: userDB.Role,
	}, nil
}

func (s *authService) Notifier(ctx context.Context, user models.CreateUserRequest) error {
	return s.notifier.SendNotification(user.Email, "Welcome", fmt.Sprintf("Hi %s Welome to the Booking Application", user.Name))
}
