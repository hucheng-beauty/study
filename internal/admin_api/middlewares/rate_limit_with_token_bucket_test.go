package middlewares

import (
	"testing"

	"github.com/gin-gonic/gin"
)

const maxConn = 100

func TestMainer(t *testing.T) {
	r := gin.Default()
	r.Use(LimitHandler(maxConn))

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	r.Run(":8080")
}
