package handler

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"xj_web_server/module/db/db"
	msg "xj_web_server/module/db/proto"
	loginProto "xj_web_server/module/login/proto"
	"xj_web_server/util"
)

type LoginDBGRPCService struct{}

//微信登陆
func (self *LoginDBGRPCService) WechatLogin(ctx context.Context, req *msg.LoginDB_C_Wechat, res *msg.LoginDB_S_Resp) error {
	login_S_Success := new(loginProto.Login_S_Success)
	rows, err := db.GetDB().DB().Query("call LSP_WechatLogin(?, ?, ?, ?, ?, ?, ?, ?)", req.AgentID, req.UserUin, req.Gender, req.NikeName, req.HeadImageUrl, req.MachineID, req.ClientAddr, req.DeviceType)
	if err != nil {
		return err
	}

	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				util.Logger.Errorf("rows 关闭错误 出错 err: %s ", err.Error())
			}
		}
	}()
	rows.Next()
	rows.Scan(&res.ErrorCode, &res.ErrorMsg, login_S_Success)
	if err != nil {
		return err
	}

	if res.ErrorCode != 0 {
		//登陆失败
		login_S_Fail := new(loginProto.Login_S_Fail)
		login_S_Fail.ErrorCode = res.ErrorCode
		login_S_Fail.ErrorMsg = res.ErrorMsg
		data, err := proto.Marshal(login_S_Fail)
		if err != nil {
			return err
		}
		copy(res.Data, data)
		return nil
	}

	data, err := proto.Marshal(login_S_Success)
	copy(res.Data, data)

	return nil
}

//手机登陆
func (self *LoginDBGRPCService) MobileLogin(ctx context.Context, req *msg.LoginDB_C_Mobile, res *msg.LoginDB_S_Resp) error {
	login_S_Success := new(loginProto.Login_S_Success)
	err := db.GetDB().DB().QueryRow("call LSP_MobileLogin(?, ?, ?, ?, ?)", req.PhoneNumber, req.Password, req.MachineID, req.DeviceType, req.ClientAddr).Scan(
		&res.ErrorCode, &res.ErrorMsg, login_S_Success)
	if err != nil {
		return err
	}

	if res.ErrorCode != 0 {
		//登陆失败
		login_S_Fail := new(loginProto.Login_S_Fail)
		login_S_Fail.ErrorCode = res.ErrorCode
		login_S_Fail.ErrorMsg = res.ErrorMsg
		data, err := proto.Marshal(login_S_Fail)
		if err != nil {
			return err
		}
		copy(res.Data, data)
		return nil
	}

	data, err := proto.Marshal(login_S_Success)
	copy(res.Data, data)

	return nil
}

//游客登陆
func (self *LoginDBGRPCService) VisitorLogin(ctx context.Context, req *msg.LoginDB_C_Visitor, res *msg.LoginDB_S_Resp) error {
	login_S_Success := new(loginProto.Login_S_Success)
	err := db.GetDB().DB().QueryRow("call LSP_VisitorLogin(?, ?, ?)", req.MachineID, req.DeviceType, req.ClientAddr).Scan(
		&res.ErrorCode, &res.ErrorMsg, login_S_Success)
	if err != nil {
		return err
	}

	if res.ErrorCode != 0 {
		//登陆失败
		login_S_Fail := new(loginProto.Login_S_Fail)
		login_S_Fail.ErrorCode = res.ErrorCode
		login_S_Fail.ErrorMsg = res.ErrorMsg
		data, err := proto.Marshal(login_S_Fail)
		if err != nil {
			return err
		}
		copy(res.Data, data)
		return nil
	}

	data, err := proto.Marshal(login_S_Success)
	copy(res.Data, data)

	return nil
}
