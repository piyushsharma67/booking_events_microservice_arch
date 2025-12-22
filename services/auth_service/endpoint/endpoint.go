package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/piyushsharma67/movie_booking/services/auth_service/service"
)

type Endpoints struct {
	SignUp   endpoint.Endpoint
	Login    endpoint.Endpoint
	Validate endpoint.Endpoint
}

func MakeEndpoints(svc service.AuthService) Endpoints {
	return Endpoints{
		SignUp: MakeSignUpEndpoint(svc),
		Login:  MakeLoginEndpoint(svc),
	}
}
