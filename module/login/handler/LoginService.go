package handler

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro"
	dbProto "xj_web_server/module/db/proto"
	msg "xj_web_server/module/login/proto"
)

const(
	LOGIN_SUCCESS = iota	//登陆成功
	LOGIN_FAIL				//登陆失败
)

var loginService LoginService

func init() {
	service := micro.NewService()

	// 初始化， 解析命令行参数等
	service.Init()
	cli := service.Client()
	// 初始化一个account服务的客户端
	loginService.loginDBCli = dbProto.NewLoginDBGRPCService("xj_web_server.service.db", cli)
}

type LoginService struct {
	loginDBCli dbProto.LoginDBGRPCService
}

//微信登陆
func (self *LoginService) wechatLogin(req *msg.ReqLogin, res *msg.RespLogin) error {
	login_C_Wechat := new(msg.Login_C_Wechat)
	err := proto.Unmarshal(req.Data, login_C_Wechat)
	if err != nil {
		return err
	}

	loginDB_C_Wechat := new(dbProto.LoginDB_C_Wechat)
	loginDB_C_Wechat.AgentID = login_C_Wechat.AgentID
	loginDB_C_Wechat.UserUin = login_C_Wechat.UserUin
	loginDB_C_Wechat.Gender = login_C_Wechat.Gender
	loginDB_C_Wechat.NikeName = login_C_Wechat.NikeName
	loginDB_C_Wechat.HeadImageUrl = login_C_Wechat.HeadImageUrl
	loginDB_C_Wechat.MachineID = login_C_Wechat.MachineID
	loginDB_C_Wechat.DeviceType = login_C_Wechat.DeviceType
	loginDB_C_Wechat.ClientAddr = req.ClientAddr
	loginDB_S_Resp, err := self.loginDBCli.WechatLogin(context.TODO(), loginDB_C_Wechat)
	if err != nil {
		return err
	}
	if loginDB_S_Resp.ErrorCode != 0 {
		//登陆失败
		res.RouteID = LOGIN_FAIL
		res.IsClose = true
	} else {
		//登陆成功
		res.RouteID = LOGIN_SUCCESS
		res.IsClose = false
	}
	copy(res.Data, loginDB_S_Resp.Data)

	return nil
}

//手机登陆
func (self *LoginService) mobileLogin(req *msg.ReqLogin, res *msg.RespLogin) error {
	login_C_Mobile := new(msg.Login_C_Mobile)
	err := proto.Unmarshal(req.Data, login_C_Mobile)
	if err != nil {
		return err
	}

	loginDB_C_Mobile := new(dbProto.LoginDB_C_Mobile)
	loginDB_C_Mobile.PhoneNumber = login_C_Mobile.PhoneNumber
	loginDB_C_Mobile.Password = login_C_Mobile.Password
	loginDB_C_Mobile.MachineID = login_C_Mobile.MachineID
	loginDB_C_Mobile.DeviceType = login_C_Mobile.DeviceType
	loginDB_C_Mobile.ClientAddr = req.ClientAddr
	loginDB_S_Resp, err := self.loginDBCli.MobileLogin(context.TODO(), loginDB_C_Mobile)
	if err != nil {
		return err
	}
	if loginDB_S_Resp.ErrorCode != 0 {
		//登陆失败
		res.RouteID = LOGIN_FAIL
		res.IsClose = true
	} else {
		//登陆成功
		res.RouteID = LOGIN_SUCCESS
		res.IsClose = false
	}
	copy(res.Data, loginDB_S_Resp.Data)

	return nil
}

//游客登陆
func (self *LoginService) visitorLogin(req *msg.ReqLogin, res *msg.RespLogin) error {
	login_C_Visitor := new(msg.Login_C_Visitor)
	err := proto.Unmarshal(req.Data, login_C_Visitor)
	if err != nil {
		return err
	}

	loginDB_C_Visitor := new(dbProto.LoginDB_C_Visitor)
	loginDB_C_Visitor.MachineID = login_C_Visitor.MachineID
	loginDB_C_Visitor.DeviceType = login_C_Visitor.DeviceType
	loginDB_C_Visitor.ClientAddr = req.ClientAddr
	loginDB_S_Resp, err := self.loginDBCli.VisitorLogin(context.TODO(), loginDB_C_Visitor)
	if err != nil {
		return err
	}
	if loginDB_S_Resp.ErrorCode != 0 {
		//登陆失败
		res.RouteID = LOGIN_FAIL
		res.IsClose = true
	} else {
		//登陆成功
		res.RouteID = LOGIN_SUCCESS
		res.IsClose = false
	}
	copy(res.Data, loginDB_S_Resp.Data)

	return nil
}
