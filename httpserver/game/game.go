package game

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xj_web_server/db"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/model"
	"xj_web_server/module"
	"xj_web_server/util"
)

type GetAllGameVersionReq struct {
	servermiddleware.BaseReq
	Platform string `json:"platform"  binding:"required"`
}

func GetAllGameVersion(c *gin.Context) {
	var resReq GetAllGameVersionReq

	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("GetAllGameVersion 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	versions, err := model.GetAllGameVersion(db.GetDBPlatform(), resReq.Platform)

	if err != nil {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: versions})

	//c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: &struct {
	//	Bundles     []model.GameVersion `json:"bundles"`
	//	GameVersion string              `json:"game_version"`
	//	GameUrl     string              `json:"game_url"`
	//}{
	//	Bundles:     versions,
	//	GameVersion: "11",
	//	GameUrl:     "22"}})

}

func GetAllGameVersionV2(c *gin.Context) {
	var resReq GetAllGameVersionReq

	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("GetAllGameVersion 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	versions, err := model.GetAllGameVersion(db.GetDBPlatform(), resReq.Platform)

	if err != nil {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	var gameUrl, downUrl string
	var versionCode string
	if resReq.Platform == "Android" {
		conf, _ := model.GetVersionConf(db.GetDBPlatform(), "GameAndroidConfig")
		if len(conf) > 0 && conf[0].Field1 != "" && conf[0].Field2 != "" {
			gameUrl = conf[0].Field1
			versionCode = conf[0].Field2
			downUrl = conf[0].Field3
		}
	} else if resReq.Platform == "IPhonePlayer" {
		conf, _ := model.GetVersionConf(db.GetDBPlatform(), "GameIOSConfig")
		if len(conf) > 0 && conf[0].Field1 != "" && conf[0].Field2 != "" {
			gameUrl = conf[0].Field1
			versionCode = conf[0].Field2
			downUrl = conf[0].Field3
		}
	} else if resReq.Platform == "WindowsEditor" {
		gameUrl = "http://xjwl-qp.oss-accelerate.aliyuncs.com/Package/WindowsEditor/RichPlayer.exe"
		versionCode = "0.1.15"
	}

	conf, err := model.GetVersionConf(db.GetDBPlatform(), "SiteConfig")
	if err != nil || len(conf) <= 0 {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}
	//
	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: &struct {
		Bundles     []model.GameVersion `json:"bundles"`
		GameVersion string              `json:"game_version"`
		GameUrl     string              `json:"game_url"`
		DownUrl     string              `json:"down_url"`
		GuestUrl    string              `json:"guest_url"`
	}{
		Bundles:     versions,
		GameVersion: versionCode,
		GameUrl:     gameUrl,
		DownUrl:     downUrl,
		GuestUrl:    conf[0].Field4,
	}})

}

type SetGameVersionReq struct {
	PokerName string `json:"poker_name"  binding:"required"`
	Hash      string `json:"hash"  binding:"required"`
	Size      int64  `json:"size"  binding:"required"`
	Platform  string `json:"platform"  binding:"required"`
}

func SetGameVersion(c *gin.Context) {
	var resReq SetGameVersionReq

	err := c.ShouldBindJSON(&resReq)
	if err != nil {
		util.Logger.Errorf("SetGameVersion 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	version, err := model.SetGameVersion(db.GetDBPlatform(), resReq.Size, resReq.PokerName, resReq.Hash, resReq.Platform)

	if err != nil {
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorSqlCode, ErrorMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: version})
}
