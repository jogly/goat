package ginfx

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/banditml/goat/envfx"
	"github.com/banditml/goat/middleware"
)

var Module = fx.Provide(func(zap *zap.Logger, env *envfx.Env) *gin.Engine {
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

type bodySpyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodySpyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ginLog(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		bodySpy := bodySpyWriter{
			ResponseWriter: c.Writer,
			body:           new(bytes.Buffer),
		}
		c.Writer = bodySpy
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		end = end.UTC()

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(time.RFC3339)),
				zap.Any("headers", c.Request.Header),
				zap.Any("response", bodySpy.body.String()),
				zap.Duration("latency", latency),
			)
		}
	}
}
