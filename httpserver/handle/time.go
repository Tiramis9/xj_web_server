package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xj_web_server/module"
	"xj_web_server/util"
	"time"
)

func CheckTime(c *gin.Context) {

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: &struct {
		TimeNowUnix int64 `json:"time_now_unix"`
	}{
		TimeNowUnix: time.Now().UnixNano() / 1000,
	},})
}
