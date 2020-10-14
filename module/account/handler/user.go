package handler

import (
	"context"
	"github.com/micro/go-micro"
	"xj_web_server/module/account/proto"
	dbProto "xj_web_server/module/db/proto"
	"xj_web_server/util"
)

var (
	dbCli dbProto.DbService
)

func init() {
	service := micro.NewService()

	// 初始化， 解析命令行参数等
	service.Init()
	cli := service.Client()
	// 初始化一个account服务的客户端
	dbCli = dbProto.NewDbService("xj_web_server.service.db", cli)
}

// User : 用于实现UserServiceHandler接口的对象
type User struct{}


// Signup : 处理用户注册请求
func (u *User) SignUp(ctx context.Context, req *proto.ReqSignUp, res *proto.RespSignUp) error {
	username := req.Username
	passwd := req.Password

	// 参数简单校验
	if len(username) < 3 || len(passwd) < 5 {
		res.Code = 1
		res.Message = "注册参数无效"
		return nil
	}

	// 对密码进行加盐及取Sha1值加密
	//encPasswd := util.MD5([]byte(passwd + cfg.PasswordSalt))
	// 将用户信息注册到用户表中
	//dbResp, err := dbcli.UserSignup(username, encPasswd)
	//if err == nil && dbResp.Suc {
	//
	//} else {
	//	res.Code = util.ErrorLackCode
	//	res.Message = "注册失败"
	//}
	res.Code = util.SuccessCode
	res.Message = "注册成功"
	return nil
}


