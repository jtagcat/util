package std

import "github.com/gin-gonic/gin"

// c.Cookie(), without ok var (instead, check if empty)
func Cookie(c *gin.Context, name string) string {
	njom, _ := c.Cookie(name)
	return njom
}
