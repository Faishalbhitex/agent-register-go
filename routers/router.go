// routers/router.go
package routers

import (
	"agent-register-go/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "agent-register"})
	})

	// Agent routes
	r.POST("/agents", handlers.RegisterAgent)
	r.GET("/agents", handlers.GetAllAgents)
	r.GET("/agents/:id", handlers.GetAgentByID)
	r.DELETE("/agents/:id", handlers.DeleteAgent)

	return r
}
