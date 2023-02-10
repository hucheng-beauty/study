package router

import (
	"net/http"

	"study/internal/admin_api/api"
	"study/internal/admin_api/api/enum"
	"study/internal/admin_api/api/request"
	"study/internal/admin_api/api/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Handler(c *gin.Context) {
	var commonRequest request.Common
	if err := c.ShouldBindBodyWith(&commonRequest, binding.JSON); err != nil {
		c.JSON(http.StatusOK, response.Error(enum.EmptyEnumRequest["Empty"],
			enum.RtCodeEnumRequest["Failed"], err.Error()))
		return
	}

	if commonRequest.Action == enum.EmptyEnumRequest["Empty"] {
		c.JSON(http.StatusOK, response.Error(enum.EmptyEnumRequest["Empty"],
			enum.RtCodeEnumRequest["Failed"], "Missing Action"))
		return
	}

	action, ok := api.Actions[commonRequest.Action]
	if !ok {
		c.JSON(http.StatusOK, response.Error(enum.EmptyEnumRequest["Empty"],
			enum.RtCodeEnumRequest["Failed"], "Not exist Action: "+commonRequest.Action))
		return
	}

	data := action.Do(c)

	c.JSON(http.StatusOK /*200*/, data)
}
