package news

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/model"
	"xj_web_server/module"
	"xj_web_server/util"
)

func NewsInfo(c *gin.Context) {
	var authReq servermiddleware.BaseReq

	err := c.ShouldBindJSON(&authReq)

	if err != nil {
		util.Logger.Errorf("NewsInfo 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	news, err := model.GetNewsInfo()

	if err != nil {
		util.Logger.Errorf("NewsInfo 接口  查询 sql语句 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: news})
}

type RecordPrizeInfoReq struct {
	servermiddleware.BaseReq
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
}

func RecordPrizeInfo(c *gin.Context) {
	var resReq RecordPrizeInfoReq

	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("RecordPrizeInfo 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	records, err := model.GetRecordPrizeInfo(resReq.Size, resReq.Page)
	if err != nil {
		util.Logger.Errorf("RecordPrizeInfo 接口  mysql err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: records})
}
