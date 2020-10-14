package exchange

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

type ConfigReq struct {
	servermiddleware.BaseReq
}

type ConfigResp struct {
	RetainCoin float32		`json:"retain_coin"`			//账户需要保留数目
	ExChangeRate float32	`json:"ex_change_rate"`			//手续费百分比
	WithdrawMaxAmount float32	`json:"withdraw_max_amount"`	//最大兑出数目
	WithdrawMinAmount float32	`json:"withdraw_min_amount"`	//最小兑出数目
	LastCodeAmounts float32		`json:"last_code_amounts"`		//剩余打码量
}

type DiamondExchangeReq struct {
	servermiddleware.BaseReq
	Type   int     `form:"type" json:"type" binding:"required"`
	Amount float64 `form:"amount" json:"amount" binding:"required"`
}

type RecordExchangeReq struct {
	servermiddleware.BaseReq
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
}

func Config(c *gin.Context) {
	// 查询账号是否存在 token解析出的uid
	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_LoadExchangeInfo(?)",
		uid,
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
		util.Logger.Errorf("DiamondExchange 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("Config 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}


	configResp := new(ConfigResp)
	rows.NextResultSet()
	err = rows.Scan(&configResp.RetainCoin, &configResp.ExChangeRate, &configResp.WithdrawMaxAmount, &configResp.WithdrawMinAmount, &configResp.LastCodeAmounts)
	if err != nil {
		util.Logger.Errorf("Config 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: configResp})
}

func DiamondExchange(c *gin.Context) {
	var resReq DiamondExchangeReq

	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("DiamondExchange 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}
	// 查询账号是否存在 token解析出的uid
	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_ApplyExchange(?,?,?,?)",
		uid,
		resReq.Amount,
		resReq.Type,
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
		util.Logger.Errorf("DiamondExchange 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("DiamondExchange 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "申请成功!"})
}

func DiamondExchangeRecord(c *gin.Context) {
	var resReq RecordExchangeReq

	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("DiamondExchangeRecord 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	// 查询账号是否存在 token解析出的uid
	uid, _ := strconv.Atoi(c.GetString("uid"))

	records, err := model.GetRecordExchangeByUid(uid, resReq.Size, resReq.Page)
	if err != nil {
		util.Logger.Errorf("DiamondExchangeRecord 接口  mysql err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: records})
}
