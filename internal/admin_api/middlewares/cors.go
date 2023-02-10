package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,"+
		"Authorization, Token, x-token, x-requested-with")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin,"+
		"Access-Control-Allow-Headers, Content-Type, Content-Disposition")
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method != http.MethodOptions {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusNoContent)
	}
}
