package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/movie_booking/services/auth_service/endpoint"
	"github.com/piyushsharma67/movie_booking/services/auth_service/service"
	"github.com/piyushsharma67/movie_booking/services/auth_service/transport"
)

func InitRoutes(srv service.AuthService) *gin.Engine {
	r := gin.Default()

	r.POST("/signup", transport.GinHandler(endpoint.MakeSignUpEndpoint(srv)))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	r.GET("/login", transport.GinHandler(endpoint.MakeLoginEndpoint(srv)))

	return r
}
