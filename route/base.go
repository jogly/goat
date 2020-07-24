package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	Resource() string
	Get(ctx *gin.Context)
	Post(ctx *gin.Context)
	Put(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type Base struct{}

func (b Base) Get(ctx *gin.Context) {
	ctx.String(http.StatusMethodNotAllowed, "method not allowed")
}

func (b Base) Post(ctx *gin.Context) {
	ctx.String(http.StatusMethodNotAllowed, "method not allowed")
}

func (b Base) Put(ctx *gin.Context) {
	ctx.String(http.StatusMethodNotAllowed, "method not allowed")
}

func (b Base) Delete(ctx *gin.Context) {
	ctx.String(http.StatusMethodNotAllowed, "method not allowed")
}

type get struct {
	Base

	handlerFunc gin.HandlerFunc
	resource    string
}

func (g get) Resource() string {
	return g.resource
}

func (g get) Get(c *gin.Context) {
	g.handlerFunc(c)
}

func Get(path string, h gin.HandlerFunc) HandlerInterface {
	return get{handlerFunc: h}
}
