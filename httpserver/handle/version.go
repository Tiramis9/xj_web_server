package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xj_web_server/module"
	"xj_web_server/util"
)

const version = "qp_server_v1.0.0"

func Version(context *gin.Context) {
	context.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "", Data: version})
}
