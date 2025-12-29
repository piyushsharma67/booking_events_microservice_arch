package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/events_booking/services/events_service/endpoints"
	"github.com/piyushsharma67/events_booking/services/events_service/logger"
	"github.com/piyushsharma67/events_booking/services/events_service/middleware"
	"github.com/piyushsharma67/events_booking/services/events_service/models"
	"github.com/piyushsharma67/events_booking/services/events_service/service"
	"github.com/piyushsharma67/events_booking/services/events_service/transport"
	"github.com/piyushsharma67/events_booking/services/events_service/utils"
)

type RoutesStruct struct {
	ginEngine *gin.Engine
	service   *service.EventService
}

func (r *RoutesStruct) InitialiseRoutes() *gin.Engine {

	if r.ginEngine == nil {
		r.ginEngine = gin.Default()

		return r.ginEngine
	}

	return r.ginEngine

}

func InitRoutes(service *service.EventService, logger logger.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "I am working fine!!",
		})
	})

	r.Use(func(c *gin.Context) {
		logger.Info("Incoming request path:", c.Request.URL.Path)
		c.Next()
	})

	r.Use(middleware.SetRoleAndIdFromHeader())
	organise := r.Group("organize")
	organise.Use(middleware.RoleAuthMiddleware(utils.ALLOWED_ORGANISER_ROUTE_ROLE))

	organise.POST("/create", transport.GinHandler(endpoints.GenerateEvent(service), func() interface{} { return &models.CreateEventRequest{} }, logger))

	return r
}
