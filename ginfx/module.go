package ginfx

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/banditml/goat/envfx"
	"github.com/banditml/goat/header"
	"github.com/banditml/goat/middleware"
)

var Module = fx.Provide(func(zap *zap.Logger, env *envfx.Env) *gin.Engine {
	zap = zap.Named("gin")
	r := gin.New()
	r.Use(gintrace.Middleware("goat"))
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
		path := c.Request.URL.Path
		span := tracer.StartSpan("web.request", tracer.ResourceName(path))
		defer span.Finish()
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		end = end.UTC()

		baggage := map[string]string{}
		span.Context().ForeachBaggageItem(func(k, v string) bool {
			baggage[k] = v
			return true
		})

		requestLogger := logger.With(
			zap.String("timestamp", end.Format(time.RFC3339Nano)),
			zap.Duration("duration", latency),
			zap.String("bandit.id", c.GetHeader(header.BanditID)),
			zap.Int("http.status_code", c.Writer.Status()),
			zap.String("http.method", c.Request.Method),
			zap.String("http.referer", c.GetHeader("Referer")),
			zap.String("http.url", c.Request.RequestURI),
			zap.String("http.useragent", c.GetHeader("User-Agent")),
			zap.String("network.client.ip", c.ClientIP()),
			zap.Int("network.bytes_written", c.Writer.Size()),
			zap.Any("http.headers", c.Request.Header),
			zap.String("dd.trace_id", strconv.Itoa(int(span.Context().TraceID()))),
			zap.String("dd.span_id", strconv.Itoa(int(span.Context().SpanID()))),
			zap.Any("dd.baggage", baggage),
		)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			sb := strings.Builder{}
			for _, e := range c.Errors.Errors() {
				sb.WriteString(e)
				sb.WriteString("\\n")
			}
			requestLogger.Error("request completed with errors",
				zap.String("errors", sb.String()))
		} else {
			requestLogger.Info("request completed")
		}
	}
}
