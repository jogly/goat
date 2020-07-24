package middleware

import (
	"errors"
	"net/http"

	"github.com/banditml/goat/header"
	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	if ctx.GetHeader(header.BanditId) == "" {
		ctx.AbortWithError(
			http.StatusUnauthorized,
			errors.New(header.BanditId+" is required"))
	}

	ctx.Next()
}
