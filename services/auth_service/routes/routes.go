package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/events_booking/services/auth_service/endpoint"
	"github.com/piyushsharma67/events_booking/services/auth_service/logger"
	"github.com/piyushsharma67/events_booking/services/auth_service/middlewares"
	"github.com/piyushsharma67/events_booking/services/auth_service/models"
	"github.com/piyushsharma67/events_booking/services/auth_service/service"
	"github.com/piyushsharma67/events_booking/services/auth_service/transport"
	"github.com/piyushsharma67/events_booking/services/auth_service/utils"
)

func InitRoutes(srv service.AuthService, logger logger.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.RequestIDMiddleware())

	r.POST("/signup", transport.GinHandler(endpoint.MakeSignUpEndpoint(srv,utils.TYPE_USER), func() interface{} { return &models.CreateUserRequest{} }, logger))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	r.POST("/login", transport.GinHandler(endpoint.MakeLoginEndpoint(srv), func() interface{} { return &models.LoginUserRequest{} }, logger))
	r.GET("/validate", transport.ValidateGinHandler(srv, logger))

	organisersRoutes:=r.Group("/organiser",transport.GinHandler(endpoint.MakeSignUpEndpoint(srv,utils.TYPE_ORGANISER), func() interface{} { return &models.CreateUserRequest{} }, logger))
	organisersRoutes.POST("/create",)

	return r
}
