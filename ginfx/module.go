package ginfx

import (
	"github.com/banditml/goat/envfx"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Provide(func(env *envfx.Env) *gin.Engine {
	r := gin.Default()
	if env.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	r.RedirectTrailingSlash = true
	r.HandleMethodNotAllowed = true
	r.RedirectFixedPath = true
	r.RemoveExtraSlash = true
	return r
})
