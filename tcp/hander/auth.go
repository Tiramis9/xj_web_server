package hander

import (
	"errors"
	redis "xj_web_server/cache"
	"xj_web_server/httpserver/wss/proto"
	"xj_web_server/util"
	"xj_web_server/util/jwt"
	"strconv"
)

type BaseMsg struct {
	Token string `json:"token"`
}

type ConnMsg struct {
	BaseMsg
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

// 验证token  连接合法性 命令号 0x00
func Auth(msg proto.Hall_C_Msg, uidCmd int32) (string, error) {
	// 缓存中查询token
	token, _ := redis.GetRedisDb().Get(util.RedisKeyToken + strconv.Itoa(int(uidCmd)) + ":").Result()

	if token != msg.Token {
		//断开连接
		util.Logger.Errorf("用户未登录：%v",msg.Token)
		return "", errors.New("token is err")
	}

	// 解析token 得到uid
	_, uid, err := jwt.EasyToken{}.ValidateToken(token)
	if err != nil {
		//断开连接
		util.Logger.Errorf("token 失效：%v", err)
		return "", err
	}

	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		//断开连接
		util.Logger.Errorf("uid 非法：%v", err)
		return "", err
	}
	if int(uidCmd) != uidInt {
		//断开连接
		util.Logger.Errorf("uid 非法：%v", err)
		return "", errors.New("uid is err")
	}
	return token, nil
}
