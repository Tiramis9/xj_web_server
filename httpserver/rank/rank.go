package rank

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"xj_web_server/db"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/model"
	"xj_web_server/module"
	"xj_web_server/util"
	"strconv"
)

type RanksResp struct {
	servermiddleware.BaseReq
	RankType  int    `form:"type" json:"type"   binding:"required"`
	Page      int    `form:"page" json:"page"   binding:"required"`
	Size      int    `form:"size" json:"size"   binding:"required"`
	BeginTime string `form:"begin_time" json:"begin_time"`
	EndTime   string `form:"end_time" json:"end_time"`
}

type RanksReq struct {
	Ranks    []model.Rank `json:"ranks"`
	UserRank model.Rank   `json:"user_rank"`
}

func RanksList(c *gin.Context) {

	var baseReq RanksResp

	err := c.ShouldBindJSON(&baseReq)

	if err != nil {
		util.Logger.Errorf("RankList 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))
	var rows *sql.Rows

	if baseReq.RankType == 1 {
		rows, err = db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserScoreRank(?,?,?)",
			baseReq.Page,
			baseReq.Size,
			uid,
		)

	} else if baseReq.RankType == 2 {

		if baseReq.BeginTime == "" || baseReq.EndTime == "" {
			c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: "参数不完整!"})
			return
		}

		rows, err = db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserProfitRank(?,?,?,?,?)",
			baseReq.BeginTime,
			baseReq.EndTime,
			baseReq.Page,
			baseReq.Size,
			uid,
		)

	} else {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: "查无此type!"})
		return
	}

	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				util.Logger.Errorf("rows 关闭错误 出错 err: %s ", err.Error())
			}
		}
	}()

	if err != nil {
		util.Logger.Errorf("RanksList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("RanksList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data RanksReq

	//查询排行榜
	var rankList = make([]model.Rank, 0)

	rows.NextResultSet()

	//i := baseReq.Size * (baseReq.Page - 1)
	var index = 0
	for rows.Next() {
		var rank model.Rank

		err = rows.Scan(&rank.UserID, &rank.GoldCoin, &rank.Nickname, &rank.FaceID, &rank.HeadImageUrl, &rank.LevelNum, &rank.RoleID, &rank.SuitID, &rank.PhotoFrameID)
		if err != nil {
			break
		}
		index++
		rank.Rank = index
		//i--

		rankList = append(rankList, rank)
	}

	//for k, _ := range rankList {
	//
	//	rankList[k].Rank = i
	//
	//}

	//查询排行榜
	//var ranks = make([]model.Rank, 0)

	//for n := len(rankList); n > 0; n-- {
	//	ranks = append(ranks, rankList[n-1])
	//}

	if err != nil {
		util.Logger.Errorf("RanksList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	//sort.Slice(rankList, func(i, j int) bool {
	//	return true
	//})

	rows.NextResultSet()

	//查询排行榜 自己的榜
	var userRank model.Rank
	if rows.Next() {
		err = rows.Scan(&userRank.UserID, &userRank.GoldCoin, &userRank.Rank, &userRank.Nickname, &userRank.FaceID, &userRank.HeadImageUrl, &userRank.LevelNum, &userRank.RoleID, &userRank.SuitID, &userRank.PhotoFrameID)

		if err != nil {
			util.Logger.Errorf("RanksList 接口  查询存储过程 err: %s ", err.Error())
			c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
			return
		}

		data.UserRank = userRank
	}

	data.Ranks = rankList

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})

}
