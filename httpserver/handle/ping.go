package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xj_web_server/module"
	"xj_web_server/util"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "pong",})
}
