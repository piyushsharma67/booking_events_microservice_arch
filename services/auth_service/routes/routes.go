package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/movie_booking/services/auth_service/endpoint"
	"github.com/piyushsharma67/movie_booking/services/auth_service/middlewares"
	"github.com/piyushsharma67/movie_booking/services/auth_service/models"
	"github.com/piyushsharma67/movie_booking/services/auth_service/service"
	"github.com/piyushsharma67/movie_booking/services/auth_service/transport"
)

func InitRoutes(srv service.AuthService) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.RequestIDMiddleware())

	r.POST("/signup", transport.GinHandler(endpoint.MakeSignUpEndpoint(srv),func() interface{} { return &models.User{} }))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	r.POST("/login", transport.GinHandler(endpoint.MakeLoginEndpoint(srv), func() interface{} { return &models.User{} }))
	r.GET("/validate", transport.ValidateGinHandler(srv))

	return r
}
