package activity

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/model"
	"xj_web_server/module"
	"xj_web_server/util"
)

type ActivityReq struct {
	servermiddleware.BaseReq
	Number int `form:"number" json:"number"` //分页数
	Page   int `from:"page" json:"page"`     //当前页
}

type ActivityResp struct {
	Total       int              `json:"total"`        //总数
	CurrentPage int              `json:"current_page"` //当前页
	List        []model.Activity `json:"list"`         //活动列表
}

func ActivityList(c *gin.Context) {

	var baseReq ActivityReq

	err := c.ShouldBindJSON(&baseReq)

	if err != nil {
		util.Logger.Errorf("ActivityList 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	// todo: 获取活动信息 分页
	var resp ActivityResp

	resp.List, resp.Total, resp.CurrentPage, err = model.GetActivityByPaging(baseReq.Number, baseReq.Page)

	if err != nil {
		util.Logger.Errorf("ActivityList 接口  sql查询 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: resp})

}
