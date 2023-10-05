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

func HandlerWithErr(f func(*gin.Context) (status int, err string)) func(*gin.Context) {
	return func(c *gin.Context) {
		code, err := f(c)
		if err != "" {
			c.HTML(code, ErrorPage, gin.H{
				"err": fmt.Sprintf("%d %s: %s", code, http.StatusText(code), err),
			})
			c.Abort()
			return
		}

		if code != 0 {
			c.Status(code)
		}
	}
}

// Wrapper for HandlerWithErr format
func HTML(c *gin.Context, status int, name string, obj any) (int, string) {
	c.HTML(status, name, obj)
	return 0, ""
}

// Wrapper HandlerWithErr format
func Redirect(c *gin.Context, status int, location string) (int, string) {
	c.Redirect(status, location)
	return 0, ""
}
