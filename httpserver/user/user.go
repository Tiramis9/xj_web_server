package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	redis "xj_web_server/cache"
	"xj_web_server/db"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/model"
	"xj_web_server/module"
	"xj_web_server/util"
	"xj_web_server/util/ztime"
	"strconv"
	"strings"
	"time"
)

var codeKeyByType map[string]CodeRedis

type CodeRedis struct {
	//验证码缓存名称
	CodeName string
	// 验证码缓存次数名称
	CodeNumberName string
	// 验证码长度
	CodeLen int
	// 验证码锁定次数
	CodeCount int
}

func init() {
	codeKeyByType = map[string]CodeRedis{
		//注册验证码
		"CODE_1": {
			CodeName:       util.RedisKeyRegisteredCode,
			CodeNumberName: util.RedisKeyRegisteredCodeNumber,
			CodeLen:        4,
			CodeCount:      5,
		},
		//忘记密码验证码
		"CODE_2": {
			CodeName:       util.RedisKeyForgotCode,
			CodeNumberName: util.RedisKeyForgotCodeNumber,
			CodeLen:        4,
			CodeCount:      5,
		},
		//绑定手机验证码
		"CODE_3": {
			CodeName:       util.RedisKeyBindingCode,
			CodeNumberName: util.RedisKeyBindingCodeNumber,
			CodeLen:        4,
			CodeCount:      5,
		},
	}
}

//代理id,账号,手机号码,密码,ip,设备类型,机器序列号
type RegisteredReq struct {
	servermiddleware.BaseReq
	Account  string `form:"account" json:"account" binding:"required"`
	Mobile   string `form:"mobile" json:"mobile" `
	Code     string `form:"code" json:"code" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	AgentID  int    `form:"agent_id" json:"agent_id"`
}

type CodeReq struct {
	servermiddleware.BaseReq
	Account string `form:"account" json:"account" binding:"required"`
	Type    int    `form:"type" json:"type" binding:"required"`
}

type ForgotPwdReq struct {
	servermiddleware.BaseReq
	Account string `form:"account" json:"account" binding:"required"`
	Code    string `form:"code" json:"code" binding:"required"`
	NewPw   string `form:"new_password" json:"new_password" binding:"required"`
}

type UpdateFaceReq struct {
	servermiddleware.BaseReq
	Type         int    `form:"type" json:"type" binding:"required"`
	FaceId       int    `form:"face_id" json:"face_id"`
	Nickname     string `form:"nickname" json:"nickname"`
	RoleID       int    `form:"role_id" json:"role_id"`               // 角色标识
	SuitID       int    `form:"suit_id" json:"suit_id"`               // 套装标识
	PhotoFrameID int    `form:"photo_frame_id" json:"photo_frame_id"` // 头像框标识
}

//	SELECT 0 AS errorCode, '' AS errorMsg,
//	paraUserID AS UserID,
//	paraGoldCoin AS UserGold,
//	paraDiamond  AS UserDiamond,
//	strPhoneNumber AS PhoneNumber,
//	paraLevelNum AS MemberOrder;

type RegisteredRsp struct {
	UserID      int64   `json:"user_id"`      //uid
	UserGold    float64 `json:"user_gold"`    //金币
	UserDiamond int64   `json:"user_diamond"` //砖石
	PhoneNumber string  `json:"phone_number"` //手机号码
	MemberOrder int64   `json:"member_order"` //等级
}

type RecordResp struct {
	servermiddleware.BaseReq
	Page      int    `form:"page" json:"page"   binding:"required"`
	Size      int    `form:"size" json:"size"   binding:"required"`
	BeginTime string `form:"begin_time" json:"begin_time"`
	EndTime   string `form:"end_time" json:"end_time"`
}

type RecordReq struct {
	RecordList []model.Record `json:"record_list"`
	TotalCount int            `json:"total_count"`
	PageCount  int            `json:"page_count"`
}

type BaseEnsuresReq struct {
	BaseEnsures    []model.Ensure `json:"base_ensures"`
	TodayCheckined int            `json:"today_checkined"`
}

type TaskBaseEnsureReq struct {
	CurrScore int `json:"curr_score"`
}

type BindingMobileResp struct {
	servermiddleware.BaseReq
	Mobile string `form:"mobile" json:"mobile"   binding:"required"`
	Code   string `form:"code" json:"code"   binding:"required"`
}

type BindingBankcardResp struct {
	servermiddleware.BaseReq
	UserName string `form:"user_name" json:"user_name"   binding:"required"`
	CardNo   string `form:"card_no" json:"card_no"   binding:"required"`
	BankName string `form:"bank_name" json:"bank_name"   binding:"required"`
}

//注册
func Registered(c *gin.Context) {
	var resReq RegisteredReq

	err := c.ShouldBindJSON(&resReq)

	if err != nil {
		util.Logger.Errorf("Registered 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	redisCode, _ := redis.GetRedisDb().Get(util.RedisKeyRegisteredCode + resReq.Mobile).Result()
	if resReq.Code != redisCode {
		util.Logger.Errorf("Registered 接口  code is %s err: %s ", redisCode, errors.New("code is err"))
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: "验证码错误"})
		return
	}
	//执行注册的存储过程
	//(IN `strAccounts` varchar(32),IN `strPhoneNumber` varchar(11),IN
	// `strLogonPass` varchar(32),IN `numAgentID`
	// int,IN `numRegisterOrigin`
	// int,IN `strMachineID`
	// varchar(32),IN `strClientIP` varchar(15),OUT `errorCode` int,OUT `errorDescribe` varchar(127))
	rows, err := db.GetDB().DB().DB.Query("CALL WSP_PW_RegisterAccounts(?,?,?,?,?,?,?)",
		resReq.Account,
		resReq.Mobile,
		resReq.Password,
		resReq.AgentID,
		resReq.NumRegisterOrigin,
		resReq.MachineID,
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
		util.Logger.Errorf("Registered 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()
	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("Registered 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}
	//查询映射到结构体
	var rsq RegisteredRsp
	rows.NextResultSet()
	rows.Next()
	//	SELECT 0 AS errorCode, '' AS errorMsg,
	//	paraUserID AS UserID,
	//	paraGoldCoin AS UserGold,
	//	paraDiamond  AS UserDiamond,
	//	strPhoneNumber AS PhoneNumber,
	//	paraLevelNum AS MemberOrder;
	err = rows.Scan(&rsq.UserID,
		&rsq.UserGold,
		&rsq.UserDiamond,
		&rsq.PhoneNumber,
		&rsq.MemberOrder)
	if err != nil {
		util.Logger.Errorf("Registered 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: rsq})
}

//忘记密码
func ForgotPwd(c *gin.Context) {
	// TODO 密码未加密
	var resReq ForgotPwdReq
	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("ForgotPwd 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}
	// 查询账号是否存在 根据 Account
	userMsg, isCheck, err := model.GetUserByAccount(resReq.Account)

	if !isCheck {
		util.Logger.Errorf("ForgotPwd 接口 account 失败 %n  %s", resReq.Account)
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: "查无此账号！"})
		return
	}

	if err != nil {
		util.Logger.Errorf("ForgotPwd 接口 account 失败 %n  %s", resReq.Account, err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: "查询账号失败！"})
		return
	}

	// 验证code是否合法
	exists := redis.GetRedisDb().Exists(util.RedisKeyForgotCode + resReq.Account).Val()

	if exists == 0 {
		util.Logger.Errorf("ForgotPwd 接口 验证码过期 %s ", resReq.Account)
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: "验证码已过期！"})
		return
	}

	code, _ := redis.GetRedisDb().Get(util.RedisKeyForgotCode + resReq.Account).Result()

	if code != resReq.Code {
		util.Logger.Errorf("ForgotPwd 接口 验证码错误 %s ", resReq.Account)
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: "验证码错误！"})
		return
	}

	// 修改密码
	pw := util.MD5(resReq.NewPw + userMsg.CodeKey)

	err = model.SetUserPasswordByUid(userMsg.Uid, pw)
	if err != nil {
		util.Logger.Errorf("ForgotPwd 接口 修改密码失败 失败 %s %s", resReq.Account, err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: "修改密码失败！"})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "修改成功！"})
}

// 修改用户信息
func UpdateInformation(c *gin.Context) {
	var resReq UpdateFaceReq
	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("UpdateUser 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}
	// 查询账号是否存在 token解析出的uid
	uid, _ := strconv.Atoi(c.GetString("uid"))
	//  根据类型 修改用户信息
	switch resReq.Type {
	case 1:
		if resReq.FaceId == 0 {
			err = errors.New("请选择头像")
		} else {
			err = model.UpdateUserFaceIdByUid(uid, resReq.FaceId)
		}
	case 2:
		if resReq.Nickname == "" {
			err = errors.New("请输入用户昵称")
		} else {
			//查询是否是敏感词
			isSen, _ := model.IsSensitiveWords(resReq.Nickname)
			if isSen {
				err = errors.New("用户昵称为敏感词,请重新输入")
			} else {
				err = model.UpdateUserNicknameByUid(uid, resReq.Nickname)
			}
		}
	case 3:
		if resReq.PhotoFrameID == 0 {
			err = errors.New("请选择头像框")
		} else {
			err = model.UpdateUserImageFrameByUid(uid, resReq.PhotoFrameID)
		}
	case 4:
		if resReq.RoleID == 0 || resReq.SuitID == 0 {
			err = errors.New("请选择装扮")
		} else {
			err = model.UpdateUserImageRoleByUid(uid, resReq.RoleID, resReq.SuitID)
		}
	}

	if err != nil {
		util.Logger.Errorf("UpdateUser 接口 UpdateFace 失败 : %s", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "修改成功!"})

}

//获取验证码
func GetCode(c *gin.Context) {
	var codeReq CodeReq

	err := c.ShouldBindJSON(&codeReq)

	if err != nil {
		util.Logger.Errorf("GetCode 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}
	codeKey, ok := codeKeyByType["CODE_"+strconv.Itoa(codeReq.Type)]
	if !ok {
		util.Logger.Errorf("GetCode 接口  参数绑定 出错 err: %s ", errors.New("type is err"))
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: errors.New("type is err").Error()})
		return
	}

	err = redis.GetRedisDb().Incr(codeKey.CodeNumberName + codeReq.Account).Err()
	//设置超时时间
	//距离凌晨多久
	day := ztime.EndOfDay(time.Now()).Sub(time.Now())
	redis.GetRedisDb().Expire(codeKey.CodeNumberName+codeReq.Account, day)
	if err != nil {
		util.Logger.Errorf("GetCode 接口  rides自增 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
		return
	}

	i, err := redis.GetRedisDb().Get(codeKey.CodeNumberName + codeReq.Account).Int()

	if err != nil {
		util.Logger.Errorf("GetCode 接口  rides获取自增 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
		return
	}

	//次数
	if i > codeKey.CodeCount {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: "当日获取验证码已达到上限！"})
		return
	}

	//发送验证码
	//vCode := "0000"
	vCode := util.GetRandDigit(codeKey.CodeLen) // 获取4位有效的验证码
	switch codeReq.Type {
	case 1:
		util.SendSms(util.REGISTERED, codeReq.Account, vCode)
	case 2:
		util.SendSms(util.ChangePassword, codeReq.Account, vCode)
	case 3:
		util.SendSms(util.BingPhoneNumber, codeReq.Account, vCode)

	}
	util.Logger.Infof("验证码%s：%s", codeReq.Account, vCode)

	err = redis.GetRedisDb().Set(codeKey.CodeName+codeReq.Account, vCode, 5*time.Minute).Err()
	if err != nil {
		util.Logger.Errorf("GetCode 接口  rides设置验证码 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "发送成功！"})
	return
}

//获取用户详情
func GetUserMsg(c *gin.Context) {

	//uid := c.GetString("uid")
	//
	//fmt.Printf(uid)
	//
	//user, err := model.GetUserById(uid)
	////var user model.User
	////err = db.GetDB().Where(` id = ? `, uid).Find(&user).Error
	//
	//if err != nil && err != gorm.ErrRecordNotFound {
	//	util.Logger.Errorf("GetUserMsg  sql 出错 err: %s ", err.Error())
	//	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
	//	return
	//}
	//
	//if err == gorm.ErrRecordNotFound {
	//	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: "登录超时！"})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: user})
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: ""})
	return
}

//获取战绩
func RecordList(c *gin.Context) {
	var baseReq RecordResp

	err := c.ShouldBindJSON(&baseReq)

	if err != nil {
		util.Logger.Errorf("RecordList 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_UserGameRecord(?,?,?,?,?)",
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
		util.Logger.Errorf("RecordList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("RecordList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data RecordReq

	var recordList = make([]model.Record, 0)

	rows.NextResultSet()

	for rows.Next() {
		var record model.Record

		err = rows.Scan(&record.RecordID, &record.UserID, &record.UserScore, &record.KindName,
			&record.GameName, &record.EnterTime, &record.LeaveTime)
		if err != nil {
			break
		}

		recordList = append(recordList, record)
	}

	data.RecordList = recordList

	if err != nil {
		util.Logger.Errorf("RecordList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&data.TotalCount, &data.PageCount)
	if err != nil {
		util.Logger.Errorf("RecordList 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
	return
}

func BaseEnsures(c *gin.Context) {
	var baseReq servermiddleware.BaseReq

	err := c.ShouldBindJSON(&baseReq)

	if err != nil {
		util.Logger.Errorf("BaseEnsures 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_LoadBaseEnsure(?)",
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
		util.Logger.Errorf("BaseEnsures 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string

	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("BaseEnsures 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data BaseEnsuresReq
	var ensures = make([]model.Ensure, 0)

	rows.NextResultSet()

	for rows.Next() {
		var ensure model.Ensure

		err = rows.Scan(&ensure.ID, &ensure.LossScore, &ensure.RewardScore, &ensure.BalanceLimit)
		if err != nil {
			break
		}

		ensures = append(ensures, ensure)
	}

	if err != nil {
		util.Logger.Errorf("BaseEnsures 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	rows.NextResultSet()
	rows.Next()

	err = rows.Scan(&data.TodayCheckined)

	if err != nil {
		util.Logger.Errorf("BaseEnsures 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	data.BaseEnsures = ensures

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: data})
}

// 领取救济金
func TakeBaseEnsure(c *gin.Context) {
	var baseReq servermiddleware.BaseReq

	err := c.ShouldBindJSON(&baseReq)

	if err != nil {
		util.Logger.Errorf("TakeBaseEnsure 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDBTreasure().DB().DB.Query("CALL WSP_PW_TakeBaseEnsure(?,?,?)",
		uid,
		strings.Split(c.Request.RemoteAddr, ":")[0],
		baseReq.MachineID,
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
		util.Logger.Errorf("TakeBaseEnsure 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("TakeBaseEnsure 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}

	var data TaskBaseEnsureReq

	rows.NextResultSet()

	rows.Next()

	err = rows.Scan(&data.CurrScore)
	if err != nil {
		util.Logger.Errorf("TakeBaseEnsure 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "领取成功!", Data: data})
}

func BindingMobile(c *gin.Context) {
	var resReq BindingMobileResp

	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("BindingMobile 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	redisCode, _ := redis.GetRedisDb().Get(util.RedisKeyBindingCode + resReq.Mobile).Result()
	if resReq.Code != redisCode {
		util.Logger.Errorf("BindingMobile 接口  code is %s err: %s ", redisCode, errors.New("code is err"))
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: "验证码错误"})
		return
	}

	// 查询账号是否存在 token解析出的uid
	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDB().DB().DB.Query("CALL WSP_PW_UserBingdingMobile(?,?,?,?,?)",
		uid,
		resReq.Mobile,
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
		util.Logger.Errorf("BindingMobile 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("BindingMobile 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "绑定成功!"})
}

// 绑定银行卡
func BindingBankcard(c *gin.Context) {
	var resReq BindingBankcardResp

	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("BindingBankcard 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}
	// TODO 后期加判断银行卡信息接口

	// 查询账号是否存在 token解析出的uid
	uid, _ := strconv.Atoi(c.GetString("uid"))

	rows, err := db.GetDB().DB().DB.Query("CALL WSP_PW_BinderBankerCard(?,?,?,?,?)",
		uid,
		resReq.UserName,
		resReq.CardNo,
		resReq.BankName,
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
		util.Logger.Errorf("BindingBankcard 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	var errorCode int64
	var errorMsg string
	rows.Next()

	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("BindingBankcard 接口  查询存储过程 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	if errorCode != util.SuccessCode {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: errorCode, ErrorMsg: errorMsg})
		return
	}
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, ErrorMsg: "绑定成功!"})
}
