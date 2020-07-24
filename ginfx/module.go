package ginfx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/banditml/goat/envfx"
	"github.com/banditml/goat/middleware"
)

var Module = fx.Provide(func(env *envfx.Env) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.Auth)
	if env.IsRelease() {
		gin.SetMode(gin.ReleaseMode)
	}
	r.HandleMethodNotAllowed = true
	r.RedirectFixedPath = true
	r.RemoveExtraSlash = true
	return r
})
