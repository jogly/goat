package ginfx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/banditml/goat/envfx"
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
