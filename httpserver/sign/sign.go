package sign

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xj_web_server/db"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/model"
	"xj_web_server/module"
	"xj_web_server/util"
	"strconv"
	"strings"
)

type SignsReq struct {
	SignReward           []model.SignReward           `json:"sign_reward"`
	SignRewardContinuous []model.SignRewardContinuous `json:"sign_reward_continuous"`
	SeriesDate           int                          `json:"series_date"`
	TodayCheckined       int                          `json:"today_checkined"`
}

type SigInReq struct {
	CurrScore float32 `json:"curr_score"`
}

type BigWheelRulesResp struct {
	servermiddleware.BaseReq
	Type int `form:"type" json:"type"`
}

type BigWheelRulesReq struct {
	Wheels          []model.Wheel `json:"wheels"`
	LimitCount      int           `json:"limit_count"`
	MinAmount       float64       `json:"min_amount"`
	AllCount        int           `json:"all_count"`
	AlreadyCount    int           `json:"already_count"`
	TotalWaterScore float64       `json:"total_water_score"`
}

type BigWheelTurntableReq struct {
	Wined     int     `json:"wined"`
	ItemIndex int     `json:"item_index"`
	ItemQuota int     `json:"item_quota"`
	LastScore float64 `json:"last_score"`
}

type BigWheelRecordResp struct {
	servermiddleware.BaseReq
	Type int `form:"type" json:"type"`
	Page int `form:"page"  binding:"required"`
	Size int `form:"size"  binding:"required"`
}

func SigIn(c *gin.Context) {
	var authReq servermiddleware.BaseReq

	err := c.ShouldBindJSON(&authReq)

	if err != nil {
		util.Logger.Errorf("SignIn 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_TakeSignin(?,?,?)",
		uid,
		strings.Split(c.Request.RemoteAddr, ":")[0],
		authReq.MachineID,
	)

	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				util.Logger.Errorf("SignIn rows 关闭错误 出错 err: %s ", err.Error())
			}
		}
	}()

	if err != nil {
		util.Logger.Errorf("SignIn 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data SigInReq

	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&data.CurrScore)

	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "签到成功！", Data: data})
}

//获取签到列表
func SigList(c *gin.Context) {
	var authReq servermiddleware.BaseReq

	err := c.ShouldBindJSON(&authReq)

	if err != nil {
		util.Logger.Errorf("%s接口:参数绑定出错err:%s", "SigList", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBPlatform().DB().DB.Query("CALL WSP_PW_LoadSigninConfig()")
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				util.Logger.Errorf("rows 关闭错误 出错 err: %s ", err.Error())
			}
		}
	}()
	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data SignsReq

	//	签到数据
	var signList = make([]model.SignReward, 0)

	rows.NextResultSet()

	for rows.Next() {
		var sign model.SignReward

		err = rows.Scan(&sign.DayID, &sign.RewardScore, &sign.Append)
		if err != nil {
			break
		}

		signList = append(signList, sign)
	}

	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	//	连续签到数据
	var signCList = make([]model.SignRewardContinuous, 0)

	rows.NextResultSet()

	for rows.Next() {
		var signC model.SignRewardContinuous

		err = rows.Scan(&signC.SeriesDay, &signC.RewardScore)
		if err != nil {
			break
		}

		signCList = append(signCList, signC)
	}

	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	rowCs, err := db.GetDBPlatform().DB().DB.Query("CALL WSP_PW_LoadUserSignin(?)", uid)

	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	defer func() {
		err := rowCs.Close()
		if err != nil {
			util.Logger.Errorf("rows 关闭错误 出错 err: %s ", err.Error())
		}
	}()

	rowCs.Next()

	err = rowCs.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	rowCs.NextResultSet()
	rowCs.Next()

	err = rowCs.Scan(&data.SeriesDate, &data.TodayCheckined)

	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	data.SignReward = signList
	data.SignRewardContinuous = signCList

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
}

func BigWheelRules(c *gin.Context) {
	var authReq BigWheelRulesResp

	err := c.ShouldBindJSON(&authReq)

	if err != nil {
		util.Logger.Errorf("%s接口:参数绑定出错err:%s", "SigList", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_LoadTurntableUserInfo(?,?)",
		uid,
		authReq.Type,
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
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data BigWheelRulesReq

	var wheels = make([]model.Wheel, 0)

	rows.NextResultSet()

	for rows.Next() {
		var wheel model.Wheel

		err = rows.Scan(&wheel.ItemIndex, &wheel.ItemQuota)

		if err != nil {
			break
		}

		wheels = append(wheels, wheel)
	}

	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	rows.NextResultSet()
	rows.Next()

	err = rows.Scan(&data.LimitCount, &data.MinAmount, &data.AllCount, &data.AlreadyCount, &data.TotalWaterScore)

	if err != nil {
		util.Logger.Errorf("SigList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	data.Wheels = wheels

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})

}

//大转盘结果
func BigWheelTurntable(c *gin.Context) {
	var authReq BigWheelRulesResp

	err := c.ShouldBindJSON(&authReq)

	if err != nil {
		util.Logger.Errorf("%s接口:参数绑定出错err:%s", "BigWheelTurntable", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_TurntableStart(?,?,?)",
		uid,
		authReq.Type,
		strings.Split(c.Request.RemoteAddr, ":")[0],
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
		util.Logger.Errorf("BigWheelTurntable 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("BigWheelTurntable 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data BigWheelTurntableReq

	rows.NextResultSet()
	rows.Next()

	err = rows.Scan(&data.Wined, &data.ItemIndex, &data.ItemQuota, &data.LastScore)

	if err != nil {
		util.Logger.Errorf("BigWheelTurntable 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
}

func BigWheelRecord(c *gin.Context) {
	var authReq BigWheelRecordResp

	err := c.ShouldBindJSON(&authReq)

	if err != nil {
		util.Logger.Errorf("%s接口:参数绑定出错err:%s", "BigWheelTurntable", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	records, err := model.GetWheelRecord(uid, authReq.Type, authReq.Size, authReq.Page)

	if err != nil {
		util.Logger.Errorf("BigWheelRecord 接口  sql语句 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: records})
}
