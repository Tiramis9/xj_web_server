package index

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	redis "xj_web_server/cache"
	"xj_web_server/httpserver/servermiddleware"
	"xj_web_server/model"
	"xj_web_server/module"
	"xj_web_server/util"
)

type StartRes struct {
	ImageUrl  string `json:"image_url"`
	LoginHost string `json:"login_host"`
}

func Star(c *gin.Context) {
	var authReq servermiddleware.BaseReq

	err := c.ShouldBindJSON(&authReq)

	if err != nil {
		util.Logger.Errorf("Index 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	imageUrl := "xxxxxxx"
	loginHost := "127.0.0.1:13000"

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: StartRes{ImageUrl: imageUrl, LoginHost: loginHost}})
}

//获取登录服务器列表
func GetHost(c *gin.Context) {
	var base servermiddleware.BaseReq
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("GetHost 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	len, err := redis.GetRedisDb().LLen(util.RedisKeyLoginServer).Result()

	if err != nil {
		util.Logger.Errorf("GetHost 接口  redis取数据 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
		return
	}

	//查询redis 取出相应的key
	keys, err := redis.GetRedisDb().LRange(util.RedisKeyLoginServer, 0, len).Result()
	if err != nil {
		util.Logger.Errorf("GetHost 接口  redis取数据 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
		return
	}

	var hostList []model.QpHost
	for _, key := range keys {

		exists := redis.GetRedisDb().Exists(key).Val()

		if exists == 1 {
			result, _ := redis.GetRedisDb().Get(key).Result()

			var host model.QpHost

			json.Unmarshal([]byte(result), &host)
			hostList = append(hostList, host)

		} else if exists == 0 {
			redis.GetRedisDb().LRem(util.RedisKeyLoginServer, 0, key)
		}

	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: hostList})

}



//获取大厅服务器列表
func GetHallHost(c *gin.Context) {
	var base servermiddleware.BaseReq
	err := c.ShouldBindJSON(&base)
	if err != nil {
		util.Logger.Errorf("GetHost 接口  参数绑定 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorLackCode, ErrorMsg: err.Error()})
		return
	}

	len, err := redis.GetRedisDb().LLen(util.RedisKeyHallServerList).Result()

	if err != nil {
		util.Logger.Errorf("GetHost 接口  redis取数据 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
		return
	}

	//查询redis 取出相应的key
	keys, err := redis.GetRedisDb().LRange(util.RedisKeyHallServerList, 0, len).Result()
	if err != nil {
		util.Logger.Errorf("GetHost 接口  redis取数据 出错 err: %s ", err.Error())
		c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.ErrorRidesCode, ErrorMsg: err.Error()})
		return
	}

	var hostList []model.QpHost
	for _, key := range keys {

		exists := redis.GetRedisDb().Exists(key).Val()

		if exists == 1 {
			result, _ := redis.GetRedisDb().Get(key).Result()

			var host model.QpHost

			json.Unmarshal([]byte(result), &host)
			hostList = append(hostList, host)

		} else if exists == 0 {
			redis.GetRedisDb().LRem(util.RedisKeyHallServerList, 0, key)
		}

	}

	c.JSON(http.StatusOK, module.ApiResp{ErrorNo: util.SuccessCode, Data: hostList})

}