package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/piyushsharma67/events_booking/services/auth_service/models"
	"github.com/piyushsharma67/events_booking/services/auth_service/service"
)

func MakeSignUpEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.CreateUserRequest)

		if req.Email == "" || req.Password == "" || req.Name == "" {
			return models.CreateUserRequest{}, errors.New("User Name,Email and Password are required")
		}

		user, err := svc.SignUp(ctx, *req)
		if err != nil {
			return nil, err
		}
		return user, nil // your SQLC User struct can be returned directly
	}
}

func MakeLoginEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.CreateUserRequest)
		user, err := svc.Login(ctx, *req)
		if err != nil {
			return nil, err
		}

		return user, nil // your SQLC User struct can be returned directly
	}
}

func MakeValidateEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.CreateUserRequest)
		user, err := svc.Login(ctx, *req)
		if err != nil {
			return struct {
				Err string `json:"error"`
			}{Err: err.Error()}, nil
		}
		return user, nil // your SQLC User struct can be returned directly
	}
}
