package agent

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xj_web_server/db"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/model"
	"xj_web_server/module"
	"xj_web_server/util"
	"strconv"
)

type DailyKnotsReq struct {
	RecordDates           []model.RecordDate           `json:"record_dates"`
	RecordGameDateAmounts []model.RecordGameDateAmount `json:"record_game_date_amounts"`
}

type UserGameInfoResp struct {
	servermiddleware.BaseReq
	StrNickname string `form:"search_nickname" json:"search_nickname"`
	Page        int    `form:"page" json:"page"   binding:"required"`
	Size        int    `form:"size" json:"size"   binding:"required"`
}

type UserGameInfoReq struct {
	UserInfo   []model.RecordGameUserInfo `json:"user_info"`
	TotalCount int                        `json:"total_count"`
	PageCount  int                        `json:"page_count"`
}

func DailyKnots(c *gin.Context) {
	var baseReq TeamsResp

	err := c.ShouldBindJSON(&baseReq)

	if err != nil {
		util.Logger.Errorf("DailyKnots 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserDailyKnots(?,?,?,?,?,?)",
		uid,
		baseReq.RankType,
		baseReq.BeginTime,
		baseReq.EndTime,
		baseReq.Page,
		baseReq.Size,

	)

	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				util.Logger.Errorf("rows 关闭错误 出错 err: %s ", err.Error())
			}
		}
	}()

	if err != nil {
		util.Logger.Errorf("DailyKnots 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("DailyKnots 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data DailyKnotsReq

	var records = make([]model.RecordDate, 0)

	rows.NextResultSet()

	for rows.Next() {
		var record model.RecordDate

		err = rows.Scan(&record.DateID, &record.KindID, &record.KindName, &record.TotalScore)
		if err != nil {
			break
		}

		records = append(records, record)
	}

	if err != nil {
		util.Logger.Errorf("DailyKnots 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var amounts = make([]model.RecordGameDateAmount, 0)

	rows.NextResultSet()

	for rows.Next() {
		var amount model.RecordGameDateAmount

		err = rows.Scan(&amount.DateID, &amount.ExchangeAmount, &amount.OrderAmount, &amount.PresentAmount)
		if err != nil {
			break
		}

		amounts = append(amounts, amount)
	}

	if err != nil {
		util.Logger.Errorf("DailyKnots 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	data.RecordDates = records
	data.RecordGameDateAmounts = amounts

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
}

func UserGameInfo(c *gin.Context) {
	var baseReq UserGameInfoResp

	err := c.ShouldBindJSON(&baseReq)

	if err != nil {
		util.Logger.Errorf("UserGameInfo 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDB().DB().DB.Query("CALL WSP_PW_UserGameInfoList(?,?,?,?)",
		uid,
		baseReq.StrNickname,
		baseReq.Page,
		baseReq.Size,

	)

	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				util.Logger.Errorf("rows 关闭错误 出错 err: %s ", err.Error())
			}
		}
	}()

	if err != nil {
		util.Logger.Errorf("UserGameInfo 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("UserGameInfo 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data UserGameInfoReq

	var users = make([]model.RecordGameUserInfo, 0)

	rows.NextResultSet()

	for rows.Next() {
		var user model.RecordGameUserInfo

		err = rows.Scan(&user.UserId, &user.NickName, &user.Diamond, &user.CollectDate)
		if err != nil {
			break
		}

		users = append(users, user)
	}

	if err != nil {
		util.Logger.Errorf("UserGameInfo 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&data.TotalCount, &data.PageCount)
	if err != nil {
		util.Logger.Errorf("UserGameInfo 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	data.UserInfo = users

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
	return

}
