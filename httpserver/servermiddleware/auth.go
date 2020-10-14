package servermiddleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	redis "xj_web_server/cache"
	"xj_web_server/module"
	"xj_web_server/util"
	"xj_web_server/util/jwt"
	//"strconv"
)

//type BaseAuthReq struct {
//	BaseReq
//	Uid int `form:"uid" json:"uid"  binding:"required"`
//}

//token验证
func BaseAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		var authReq BaseReq

		req, err := c.GetRawData()
		if err != nil {
			util.Logger.Errorf("BaseAuth  参数绑定 出错 err: %s ", err.Error())
			c.AbortWithStatusJSON(http.StatusOK, module.ApiResp{ErrorNo: http.StatusForbidden, ErrorMsg: err.Error()})
			return
		}
		//传递参数到下个中间件
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(req)) // 关键点

		err = json.Unmarshal(req, &authReq)

		if err != nil {
			util.Logger.Errorf("BaseAuth  参数绑定 出错 err: %s ", err.Error())
			c.AbortWithStatusJSON(http.StatusOK, module.ApiResp{ErrorNo: http.StatusForbidden, ErrorMsg: err.Error()})
			return
		}

		token := c.GetHeader("token")

		if token == "" {
			//权限异常
			c.AbortWithStatusJSON(http.StatusOK, module.ApiResp{
				ErrorNo:  http.StatusForbidden,
				ErrorMsg: http.StatusText(http.StatusForbidden),
			})
			return
		}
		et := jwt.EasyToken{}
		valid, tokenUid, err := et.ValidateToken(token)

		if !valid {
			if err != nil {
				util.Logger.Errorf("BaseAuth  token 验证 出错 err: %s\n%s ", token, err.Error())
			}
			c.AbortWithStatusJSON(http.StatusOK, module.ApiResp{
				ErrorNo:  http.StatusForbidden,
				ErrorMsg: "token validate failed, please login again.",
			})
			return
		}

		//验证token是否存在
		//1,查询redis
		tokenRedis, err := redis.GetRedisDb().Get(util.RedisKeyToken + tokenUid + ":").Result()

		if tokenRedis != token {
			if err != nil {
				util.Logger.Errorf("BaseAuth  token 验证 出错 err: 客户端:%s\nredis:%s\nerr:%s ", tokenRedis,token, err.Error())
			}else{
				util.Logger.Errorf("BaseAuth  token 验证 出错 err: 客户端:%s\nredis:%s", tokenRedis,token)
			}
			c.AbortWithStatusJSON(http.StatusOK, module.ApiResp{
				ErrorNo:  http.StatusForbidden,
				ErrorMsg: "token failed, please login again.",
			})
			return
		}

		c.Set("uid", tokenUid)

	}

}
