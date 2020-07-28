package ginfx

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/banditml/goat/envfx"
	"github.com/banditml/goat/header"
	"github.com/banditml/goat/middleware"
)

var Module = fx.Provide(func(zap *zap.Logger, env *envfx.Env) *gin.Engine {
	zap = zap.Named("gin")
	r := gin.New()
	r.Use(ginLog(zap))
	r.Use(gin.Recovery())
	r.Use(middleware.Auth)
	if env.IsRelease() {
		gin.SetMode(gin.ReleaseMode)
	}
	r.HandleMethodNotAllowed = true
	r.RedirectFixedPath = true
	r.RemoveExtraSlash = true
	return r
})

func ginLog(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		end = end.UTC()

		logger = logger.With(
			zap.String("timestamp", end.Format(time.RFC3339)),
			zap.Duration("duration", latency),
			zap.String("bandit.id", c.Request.Header.Get(header.BanditID)),
			zap.Int("http.status_code", c.Writer.Status()),
			zap.String("http.method", c.Request.Method),
			zap.String("http.url_details.path", path+"?"+query),
			zap.String("network.client.ip", c.ClientIP()),
			zap.Int("network.bytes_written", c.Writer.Size()),
			zap.Any("http.headers", c.Request.Header),
		)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			sb := strings.Builder{}
			for _, e := range c.Errors.Errors() {
				sb.WriteString(e)
				sb.WriteString("\\n")
			}
			logger.Error("errors encountered during request",
				zap.String("errors", sb.String()))
		} else {
			logger.Info(path)
		}
	}
}
