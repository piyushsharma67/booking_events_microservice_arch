package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/movie_booking/services/movies_service/databases"
)

func InitialiseRoutes(db databases.Database) *gin.Engine {

	r := gin.Default()

	r.GET("/health", HealthCheck)
	r.POST("/signup", SignUpHandler(db))
	r.GET("/login", LoginHandler(db))
	//adding the route to validate the token which would be used by ngin auth_request module
	r.GET("/validate", ValidateHandler())

	return r
}
