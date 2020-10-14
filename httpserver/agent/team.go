package agent

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	redis "xj_web_server/cache"
	"xj_web_server/db"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/model"
	"xj_web_server/module"
	"xj_web_server/util"
	"strconv"
	"strings"
	"time"
)

type OneResp struct {
	servermiddleware.BaseReq
	UserId int `form:"user_id" json:"user_id"`
	Page   int `form:"page" json:"page"   binding:"required"`
	Size   int `form:"size" json:"size"   binding:"required"`
}

type OneMsg struct {
	TotalRewardAmount float32      `json:"total_reward_amount"`
	DayRewardAmount   float32      `json:"day_reward_amount"`
	UserCount         int          `json:"user_count"`
	TotalCount        int          `json:"-"`
	PageCount         int          `json:"-"`
	Records           []OneRecords `json:"records"`
}

type OneRecords struct {
	UserID         int     `json:"user_id"`
	RechargeAmount float32 `json:"recharge_amount"`
	RewardAmount   float32 `json:"reward_amount"`
	RewardDate     string  `json:"reward_date"`
}

type TeamsResp struct {
	servermiddleware.BaseReq
	RankType  int    `form:"type" json:"type"`
	Page      int    `form:"page" json:"page"   binding:"required"`
	Size      int    `form:"size" json:"size"   binding:"required"`
	BeginTime string `form:"begin_time" json:"begin_time"`
	EndTime   string `form:"end_time" json:"end_time"`
}

type TeamsReq struct {
	Records    []model.RecordGame     `json:"records"`
	TotalCount int                    `json:"total_count"`
	PageCount  int                    `json:"page_count"`
	Amounts    model.RecordGameAmount `json:"amounts"`
}

type DiamondChangeLogResp struct {
	servermiddleware.BaseReq
	StrNickname string `form:"search_nickname" json:"search_nickname"`
	RankType    int    `form:"type" json:"type"`
	BeginTime   string `form:"begin_time" json:"begin_time" binding:"required"`
	EndTime     string `form:"end_time" json:"end_time" binding:"required"`
	Page        int    `form:"page" json:"page"   binding:"required"`
	Size        int    `form:"size" json:"size"   binding:"required"`
}

type DiamondChangeLogReq struct {
	Logs       []model.RecordDiamondChangeLog `json:"diamond_change_logs"`
	TotalCount int                            `json:"total_count"`
	PageCount  int                            `json:"page_count"`
}

type MyPromoteMsg struct {
	ParentID          int     `json:"parent_id"`
	UserID            int     `json:"user_id"`
	AllCount          int     `json:"all_count"`
	TotalCount        int     `json:"total_count"`
	DireAllCount      int     `json:"dire_all_count"`
	DireAgentAllCount int     `json:"dire_agent_all_count"`
	DireCount         int     `json:"dire_count"`
	DireAgentCount    int     `json:"dire_agent_count"`
	YesterdayAmount   float64 `json:"yesterday_amount"`
	TodayAmount       float64 `json:"today_amount"`
	TotalWaterScore   float64 `json:"total_water_score"`
	DireWaterScore    float64 `json:"dire_water_score"`
	TotalAgentAmount  float64 `json:"total_agent_amount"`
	CanAgentAmount    float64 `json:"can_agent_amount"`
	MyPromoteUrl      string  `json:"my_promote_url"`
}

type TeamDireReq struct {
	TeamDireMsg  TeamDireMsg    `json:"team_dire_msg"`
	TeamDireList []TeamDireList `json:"team_dire_list"`
}

type TeamDireMsg struct {
	UserID          int64   `json:"user_id"`
	DireCount       int64   `json:"dire_count"`
	DireAgentCount  int64   `json:"dire_agent_count"`
	TodayWaterScore float32 `json:"today_water_score"`
}

type TeamDireList struct {
	UserID               int64   `json:"user_id"`
	TodayWaterScore      float64 `json:"today_water_score"`
	TodayTotalWaterScore float64 `json:"today_total_water_score"`
	AllChildCount        int64   `json:"all_child_count"`
	DireChildCount       int64   `json:"dire_child_count"`
}

type TakeAgentRoyaltyResp struct {
	servermiddleware.BaseReq
	//DecAmount float64 `form:"dec_amount" json:"dec_amount"`
}

func TeamList(c *gin.Context) {
	var baseReq TeamsResp

	err := c.ShouldBindJSON(&baseReq)

	if err != nil {
		util.Logger.Errorf("TeamList 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserGameStatistics(?,?,?,?,?,?)",
		baseReq.RankType,
		uid,
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
		util.Logger.Errorf("TeamList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("TeamList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}
	var data TeamsReq

	var records = make([]model.RecordGame, 0)

	rows.NextResultSet()

	for rows.Next() {
		var record model.RecordGame

		err = rows.Scan(&record.UserID, &record.NickName, &record.KindID, &record.KindName, &record.TotalScore)
		if err != nil {
			break
		}

		records = append(records, record)
	}

	if err != nil {
		util.Logger.Errorf("TeamList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&data.TotalCount, &data.PageCount)
	if err != nil {
		util.Logger.Errorf("TeamList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	var amount model.RecordGameAmount

	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&amount.UserId, &amount.ExchangeAmount, &amount.OrderAmount, &amount.PresentAmount)
	if err != nil {
		util.Logger.Errorf("TeamList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	data.Records = records
	data.Amounts = amount

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})

}

func DiamondChangeLog(c *gin.Context) {
	var baseReq DiamondChangeLogResp

	err := c.ShouldBindJSON(&baseReq)

	if err != nil {
		util.Logger.Errorf("DiamondChangeLog 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserDiamondChangeLog(?,?,?,?,?,?,?)",
		uid,
		baseReq.StrNickname,
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
		util.Logger.Errorf("DiamondChangeLog 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("DiamondChangeLog 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data DiamondChangeLogReq

	var logs = make([]model.RecordDiamondChangeLog, 0)

	rows.NextResultSet()

	for rows.Next() {
		var log model.RecordDiamondChangeLog

		err = rows.Scan(&log.UserId, &log.NickName, &log.CapitalTypeID, &log.CapitalAmount, &log.LastAmount, &log.LogDate)
		if err != nil {
			break
		}

		logs = append(logs, log)
	}

	if err != nil {
		util.Logger.Errorf("TeamList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&data.TotalCount, &data.PageCount)
	if err != nil {
		util.Logger.Errorf("TeamList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	data.Logs = logs

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
}

func MyPromote(c *gin.Context) {
	var base servermiddleware.BaseReq
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("MyPromote 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserSpreaderInfo(?)",
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
		util.Logger.Errorf("MyPromote 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("MyPromote 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data MyPromoteMsg
	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&data.ParentID, &data.UserID, &data.TotalCount, &data.DireCount, &data.DireAgentCount, &data.YesterdayAmount, &data.TodayAmount, &data.TotalWaterScore, &data.DireWaterScore, &data.TotalAgentAmount, &data.CanAgentAmount, &data.AllCount, &data.DireAllCount, &data.DireAgentAllCount)

	if err != nil {
		util.Logger.Errorf("MyPromote 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	//todo 推广链接
	data.MyPromoteUrl = "https://www.xjwl.com/agent/" + strconv.Itoa(uid)

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
}

func One(c *gin.Context) {
	var base OneResp
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("One 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_RechargeRewardTotal(?)",
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
		util.Logger.Errorf("One 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("One 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data OneMsg
	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&data.TotalRewardAmount, &data.DayRewardAmount, &data.UserCount)

	if err != nil {
		util.Logger.Errorf("One 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	rows, err = db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_RechargeReward(?,?,?,?)",
		uid,
		base.UserId,
		base.Page,
		base.Size,
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
		util.Logger.Errorf("One 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("One 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	rows.NextResultSet()

	var oneRecords = make([]OneRecords, 0)
	for rows.Next() {
		var oneRecord OneRecords

		err = rows.Scan(&oneRecord.UserID, &oneRecord.RechargeAmount, &oneRecord.RewardAmount, &oneRecord.RewardDate)
		if err != nil {
			break
		}

		oneRecords = append(oneRecords, oneRecord)
	}
	if err != nil {
		util.Logger.Errorf("One 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	data.Records = oneRecords

	rows.NextResultSet()
	rows.Next()
	err = rows.Scan(&data.TotalCount, &data.PageCount)

	if err != nil {
		util.Logger.Errorf("One 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
}

func TeamDire(c *gin.Context) {
	var base TeamsResp
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("TeamDire 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	if base.Size <= 0 || base.Page <= 0 {
		util.Logger.Errorf("TeamDire 接口  参数绑定 出错 err: %v ", base)
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: "参数错误"})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	var data TeamDireReq

	// 验证code是否合法
	exists := redis.GetRedisDb().Exists(util.RedisKeyTeam + strconv.Itoa(uid) + util.RedisKeyTeamDire).Val()

	if exists != 0 {

		teamDireStr, err := redis.GetRedisDb().Get(util.RedisKeyTeam + strconv.Itoa(uid) + util.RedisKeyTeamDire).Result()

		if teamDireStr != "" {

			if err != nil {
				util.Logger.Errorf("TeamDire 接口  redis缓存转换 出错 err: %s ", err.Error())
				c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
				return
			}

			err = json.Unmarshal([]byte(teamDireStr), &data)

			if err != nil {
				util.Logger.Errorf("TeamDire 接口  redis缓存转换 出错 err: %s ", err.Error())
				c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
				return
			}

			//base.Page
			//base.Size

			if len(data.TeamDireList) < (base.Page-1)*base.Size {
				data.TeamDireList = make([]TeamDireList, 0)
			} else if len(data.TeamDireList) < base.Page*base.Size {
				data.TeamDireList = data.TeamDireList[(base.Page-1)*base.Size:]
			} else {
				data.TeamDireList = data.TeamDireList[(base.Page-1)*base.Size : (base.Page)*base.Size]
			}

			c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
			return
		}
	}

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserAgentTotalInfo(?)",
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
		util.Logger.Errorf("TeamDire 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("TeamDire 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var teamDireMsg TeamDireMsg

	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&teamDireMsg.UserID, &teamDireMsg.DireCount, &teamDireMsg.DireAgentCount, &teamDireMsg.TodayWaterScore)

	if err != nil {
		util.Logger.Errorf("TeamDire 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var teamDireList = make([]TeamDireList, 0)

	rows.NextResultSet()

	for rows.Next() {
		var teamDire TeamDireList

		err = rows.Scan(&teamDire.UserID, &teamDire.TodayWaterScore, &teamDire.TodayTotalWaterScore, &teamDire.AllChildCount, &teamDire.DireChildCount)
		if err != nil {
			break
		}

		teamDireList = append(teamDireList, teamDire)
	}

	if err != nil {
		util.Logger.Errorf("TeamDire 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	data.TeamDireMsg = teamDireMsg
	data.TeamDireList = teamDireList

	bytes, err := json.Marshal(&data)

	if err != nil {
		util.Logger.Errorf("TeamDire 接口  rides缓存失败 出错 err: %s ", err.Error())
	} else {
		err = redis.GetRedisDb().Set(util.RedisKeyTeam+strconv.Itoa(uid)+util.RedisKeyTeamDire, string(bytes), 30*time.Minute).Err()
		if err != nil {
			util.Logger.Errorf("TeamDire 接口  rides缓存失败 出错 err: %s ", err.Error())
		}
	}
	if len(teamDireList) < (base.Page-1)*base.Size {
		data.TeamDireList = make([]TeamDireList, 0)
	} else if len(teamDireList) < base.Page*base.Size {
		data.TeamDireList = teamDireList[(base.Page-1)*base.Size:]
	} else {
		data.TeamDireList = teamDireList[(base.Page-1)*base.Size : (base.Page)*base.Size]
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})

}

func GetAgentRotyaltyLevel(c *gin.Context) {
	var base servermiddleware.BaseReq
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("GetAgentRotyaltyLevel 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	levels, err := model.GetAllAgentRoyaltyLevel()
	if err != nil {
		util.Logger.Errorf("GetAgentRotyaltyLevel 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: levels})
}

type resAgentRotyaltyConfig struct {
	ReChangeValue int64  `json:"re_change_value" xorm:"-"` //充值奖励
	CanReChange   int64  `json:"can_rechange" xorm:"-"`    //返佣有效充值限额
	MinCont       int64  `json:"min_cont" xorm:"-"`        //返佣团队最低有效人数
	MyPromoteUrl  string `json:"my_promote_url"`           //推广链接
}

func GetAgentRotyaltyConfig(c *gin.Context) {
	var base servermiddleware.BaseReq
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("GetAgentRotyaltyConfig 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	levels, err := model.GetSystemStatusInfo()
	if err != nil || len(levels) < 3 {
		util.Logger.Errorf("GetAgentRotyaltyConfig 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	uid, _ := strconv.Atoi(c.GetString("uid"))
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: resAgentRotyaltyConfig{
		ReChangeValue: levels[0].StatusValue,
		CanReChange:   levels[1].StatusValue,
		MinCont:       levels[2].StatusValue,
		//todo 推广链接
		MyPromoteUrl: "https://www.xjwl.com/agent/" + strconv.Itoa(uid),
	}})
}

type TeamAchievementReq struct {
	DateTime string `json:"date_time"`
	//DireNewWaterScore  float64 `json:"dire_new_water_score"`
	DireWaterScore float64 `json:"dire_water_score"`
	//TotalNewWaterScore float64 `json:"total_new_water_score"`
	TotalWaterScore float64 `json:"total_water_score"`
	RoyaltyAmount   float64 `json:"royalty_amount"`
}

func TeamAchievement(c *gin.Context) {
	var base TeamsResp
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("TeamAchievement 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	if base.Size <= 0 || base.Page <= 0 {
		util.Logger.Errorf("TeamAchievement 接口  参数绑定 出错 err: %v ", base)
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: "参数错误"})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	var data = make([]TeamAchievementReq, 0)

	// 验证code是否合法
	exists := redis.GetRedisDb().Exists(util.RedisKeyTeam + strconv.Itoa(uid) + util.RedisKeyTeamAchievement).Val()

	if exists != 0 {
		teamDireStr, err := redis.GetRedisDb().Get(util.RedisKeyTeam + strconv.Itoa(uid) + util.RedisKeyTeamAchievement).Result()

		if teamDireStr != "" {

			if err != nil {
				util.Logger.Errorf("TeamAchievement 接口  redis缓存转换 出错 err: %s ", err.Error())
				c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
				return
			}

			err = json.Unmarshal([]byte(teamDireStr), &data)

			if err != nil {
				util.Logger.Errorf("TeamAchievement 接口  redis缓存转换 出错 err: %s ", err.Error())
				c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
				return
			}

			//base.Page
			//base.Size

			if len(data) < (base.Page-1)*base.Size {
				data = []TeamAchievementReq{}
			} else if len(data) < base.Page*base.Size {
				data = data[(base.Page-1)*base.Size:]
			} else {
				data = data[(base.Page-1)*base.Size : (base.Page)*base.Size]
			}

			c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
			return
		}
	}

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserAchievement(?)",
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
		util.Logger.Errorf("TeamAchievement 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("TeamAchievement 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	rows.NextResultSet()

	for rows.Next() {
		var ac TeamAchievementReq

		err = rows.Scan(&ac.DateTime, &ac.DireWaterScore, &ac.TotalWaterScore, &ac.RoyaltyAmount)

		data = append(data, ac)
	}

	if err != nil {
		util.Logger.Errorf("TeamAchievement 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	bytes, err := json.Marshal(&data)

	if err != nil {
		util.Logger.Errorf("TeamAchievement 接口  rides缓存失败 出错 err: %s ", err.Error())
	} else {
		err = redis.GetRedisDb().Set(util.RedisKeyTeam+strconv.Itoa(uid)+util.RedisKeyTeamAchievement, string(bytes), 30*time.Minute).Err()
		if err != nil {
			util.Logger.Errorf("TeamAchievement 接口  rides缓存失败 出错 err: %s ", err.Error())
		}
	}

	if len(data) < (base.Page-1)*base.Size {
		data = make([]TeamAchievementReq, 0)
	} else if len(data) < base.Page*base.Size {
		data = data[(base.Page-1)*base.Size:]
	} else {
		data = data[(base.Page-1)*base.Size : (base.Page)*base.Size]
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
}

type TeamAchievementFormReq struct {
	UserID          int64   `json:"user_id"`
	OwnWaterScore   float64 `json:"own_water_score"`
	DireWaterScore  float64 `json:"dire_water_score"`
	TotalWaterScore float64 `json:"total_water_score"`
}

func TeamAchievementForm(c *gin.Context) {
	var base TeamsResp
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("TeamAchievementForm 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	if base.Size <= 0 || base.Page <= 0 || base.BeginTime == "" {
		util.Logger.Errorf("TeamAchievementForm 接口  参数绑定 出错 err: %v ", base)
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: "参数错误"})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	var data []TeamAchievementFormReq

	// 验证code是否合法
	exists := redis.GetRedisDb().Exists(util.RedisKeyTeam + strconv.Itoa(uid) + util.RedisKeyTeamAchievementForm).Val()
	if exists != 0 {
		teamDireStr, err := redis.GetRedisDb().Get(util.RedisKeyTeam + strconv.Itoa(uid) + util.RedisKeyTeamAchievementForm).Result()

		if teamDireStr != "" {

			if err != nil {
				util.Logger.Errorf("TeamAchievementForm 接口  redis缓存转换 出错 err: %s ", err.Error())
				c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
				return
			}

			err = json.Unmarshal([]byte(teamDireStr), &data)

			if err != nil {
				util.Logger.Errorf("TeamAchievementForm 接口  redis缓存转换 出错 err: %s ", err.Error())
				c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
				return
			}

			//base.Page
			//base.Size

			if len(data) < (base.Page-1)*base.Size {
				data = []TeamAchievementFormReq{}
			} else if len(data) < base.Page*base.Size {
				data = data[(base.Page-1)*base.Size:]
			} else {
				data = data[(base.Page-1)*base.Size : (base.Page)*base.Size]
			}

			c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
			return
		}
	}

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserAchievementFrom(?,?)",
		uid,
		base.BeginTime,
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
		util.Logger.Errorf("TeamAchievementForm 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("TeamAchievementForm 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	rows.NextResultSet()

	for rows.Next() {
		var ac TeamAchievementFormReq

		err = rows.Scan(&ac.UserID, &ac.OwnWaterScore, &ac.DireWaterScore, &ac.TotalWaterScore)

		data = append(data, ac)
	}

	if err != nil {
		util.Logger.Errorf("TeamAchievementForm 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	bytes, err := json.Marshal(&data)

	if err != nil {
		util.Logger.Errorf("TeamAchievementForm 接口  rides缓存失败 出错 err: %s ", err.Error())
	} else {
		err = redis.GetRedisDb().Set(util.RedisKeyTeam+strconv.Itoa(uid)+util.RedisKeyTeamAchievementForm, string(bytes), 30*time.Minute).Err()
		if err != nil {
			util.Logger.Errorf("TeamAchievementForm 接口  rides缓存失败 出错 err: %s ", err.Error())
		}
	}

	if len(data) < (base.Page-1)*base.Size {
		data = []TeamAchievementFormReq{}
	} else if len(data) < base.Page*base.Size {
		data = data[(base.Page-1)*base.Size:]
	} else {
		data = data[(base.Page-1)*base.Size : (base.Page)*base.Size]
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
}

//func TeamPromote(c *gin.Context) {
//	var base servermiddleware.BaseReq
//	err := c.ShouldBindJSON(&base)
//	if err != nil {
//		util.Logger.Errorf("TeamPromote 接口  参数绑定 出错 err: %s ", err.Error())
//		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
//		return
//	}
//
//	uid, _ := strconv.Atoi(c.GetString("uid"))
//
//	//todo 生成推广链接
//	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: &struct{ Addr string }{Addr: "https://www.xjwl.com/agent/" + strconv.Itoa(uid)}})
//}

func TeamTakeAgentRoyalty(c *gin.Context) {
	var base TakeAgentRoyaltyResp
	err := c.ShouldBindJSON(&base)
	//if err != nil || base.DecAmount <= 0 {
	//	util.Logger.Errorf("TeamPromote 接口  参数绑定 出错 err: %s ", err.Error())
	//	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
	//	return
	//}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_TakeAgentRoyalty(?,?)",
		uid,
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
		util.Logger.Errorf("TeamAchievementForm 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("TeamAchievementForm 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "领取成功"})

}

func TeamTakeAgentRecord(c *gin.Context) {
	var base servermiddleware.BaseReq
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("TeamTakeAgentRecord 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	records, err := model.GetAgentRoyaltyTakeRecordByUid(uid)

	if err != nil {
		util.Logger.Errorf("GetAgentRoyaltyTakeRecordByUid mysql err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: records})
}
