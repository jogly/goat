package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(g *gin.Engine) error {
	g.GET("/health", healthCheck)
	return nil
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
