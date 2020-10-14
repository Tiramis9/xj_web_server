package handler

import (
	"context"
	msg "xj_web_server/module/login/proto"
)

const (
	WECHAT_LOGIN = iota		//微信登陆
	MOBILE_LOGIN			//手机登陆
	VISITOR_REGIST			//游客登陆
)

// LoginGRPCService : 用于实现GRPC接口的对象
type LoginGRPCService struct{
}

func (self *LoginGRPCService) Handler(ctx context.Context, req *msg.ReqLogin, res *msg.RespLogin) error {
	switch (req.RouteID) {
	case WECHAT_LOGIN:
		err := loginService.wechatLogin(req, res)
		if err != nil {
			return err
		}
	case MOBILE_LOGIN:
		err := loginService.mobileLogin(req, res)
		if err != nil {
			return err
		}
	case VISITOR_REGIST:
		err := loginService.visitorLogin(req, res)
		if err != nil {
			return err
		}
	default:
		res.RouteID = -1	//"无效指令"
	}

	return nil
}
