package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// c.Cookie(), without ok var (instead, check if empty)
func Cookie(c *gin.Context, name string) string {
	njom, _ := c.Cookie(name)
	return njom
}

// Do not change this outside of init()
var ErrorPage = "error.html"

func HandlerWithErr(f func(*gin.Context, *Context) (status int, err string)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		code, err := f(ctx, &Context{ctx})
		if err != "" {
			ctx.HTML(code, ErrorPage, gin.H{
				"err": fmt.Sprintf("%d %s: %s", code, http.StatusText(code), err),
			})
			ctx.Abort()
			return
		}

		if code != 0 {
			ctx.Status(code)
		}
	}
}

type Context struct {
	ctx *gin.Context
}

// Wrapper for HandlerWithErr format
func (w *Context) HTML(status int, name string, obj any) (int, string) {
	w.ctx.HTML(status, name, obj)
	return 0, ""
}

// Wrapper HandlerWithErr format
func (w *Context) Redirect(status int, location string) (int, string) {
	w.ctx.Redirect(status, location)
	return 0, ""
}

// Wrapper HandlerWithErr format
func (w *Context) Cookie(name string) string {
	return Cookie(w.ctx, name)
}
