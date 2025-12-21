package routes

import "github.com/gin-gonic/gin"

func HealthCheck(g *gin.Context) {
	g.JSON(200, gin.H{"status": "healthy"})
}
