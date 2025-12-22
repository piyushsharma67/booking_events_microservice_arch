package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
)

func GinHandler(e endpoint.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request interface{}

		if err:=c.ShouldBindBodyWithJSON(request);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		// Call the Go Kit endpoint
		resp, err := e(c.Request.Context(), request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}