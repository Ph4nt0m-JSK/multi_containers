package xdd

import (
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

func init() {
	server := core.Server
	server.GET(Web.Get("path", "/web"), func(c *gin.Context) {
		c.String(200, "11111111111111111111111111")
	})
}
