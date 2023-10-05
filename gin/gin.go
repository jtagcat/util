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

// Uses error.html provided by consumer
func HandlerWithErr(f func(*gin.Context) (status int, err string)) func(*gin.Context) {
	return func(c *gin.Context) {
		code, err := f(c)
		if err != "" {
			c.HTML(code, "error.html", gin.H{
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
