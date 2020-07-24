package middleware

import (
	"errors"
	"net/http"

	"github.com/banditml/goat/header"
	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	if ctx.FullPath() == "/health" || ctx.FullPath() == "" {
		ctx.Next()
		return
	}
	if ctx.GetHeader(header.BanditId) == "" {

		ctx.AbortWithError(
			http.StatusUnauthorized,
			errors.New(ctx.FullPath()+" is protected: "+header.BanditId+" is required"))
		return
	}

	ctx.Next()
}
