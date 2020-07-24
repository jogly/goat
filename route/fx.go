package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Engine   *gin.Engine
	Logger   *zap.Logger
	Handlers []HandlerInterface `group:"handlers"`
}

func Fx(constructor interface{}) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  "handlers",
			Target: constructor,
		},
	)
}

func Register(p Params) error {
	p.Engine.GET("/health", healthCheck)

	for _, h := range p.Handlers {
		p.Engine.GET("/"+h.Resource(), h.Get)
		p.Engine.GET("/"+h.Resource()+"/:uuid", h.Get)
		p.Engine.POST("/"+h.Resource(), h.Post)
		p.Engine.PUT("/"+h.Resource()+"/:uuid", h.Put)
	}

	return nil
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
