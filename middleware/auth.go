// Package middleware contains various Gin middlewares
package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/banditml/goat/header"
)

func Auth(ctx *gin.Context) {
	if ctx.FullPath() == "/health" || ctx.FullPath() == "" {
		ctx.Next()
		return
	}
	if ctx.GetHeader(header.BanditID) == "" {
		ctx.AbortWithError(
			http.StatusUnauthorized,
			errors.New(ctx.FullPath()+" is protected: "+header.BanditID+" is required"))
		return
	}

	ctx.Next()
}
