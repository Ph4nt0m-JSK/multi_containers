package xdd

import (
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

func init() {
	core.Tail = ""
	core.Server.GET("/", func(c *gin.Context) {
		c.String(200, "11111111111111111111111111")
	})
}
