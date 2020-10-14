package apigw

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"log"
	"net/http"
	"xj_web_server/module"
	userProto "xj_web_server/module/account/proto"
	"xj_web_server/module/selector"
)

var (
	userCli userProto.UserService
)

func init() {

	service := micro.NewService()

	client.DefaultClient = client.NewClient(
		client.Selector(selector.FirstNodeSelector()),
	)
	// 初始化， 解析命令行参数等
	service.Init(
		micro.AfterStart(func() error {

			return nil
		}),
		micro.BeforeStart(func() error {
			return nil
		}),
		micro.AfterStop(func() error {
			return nil
		}),
	)

	cli := service.Client()
	// 初始化一个account服务的客户端
	userCli = userProto.NewUserService("xj_web_server.service.user", cli)
}

func Call(i int, client client.Client) {
	// Create new request to service go.micro.srv.example, method Example.Call
	req := client.NewRequest("xj_web_server.service.user", "User.SignUp", &userProto.ReqSignUp{
		Username: "John",
		Password: "1333",
	})

	rsp := &userProto.RespSignUp{}

	// Call service
	if err := client.Call(context.Background(), req, rsp); err != nil {
		fmt.Println("call err: ", err, rsp)
		return
	}

	fmt.Println("Call:", i, "rsp:", rsp.Message)
}

// DoSignUpHandler : 处理注册post请求
func DoSignUpHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	resp, err := userCli.SignUp(context.TODO(), &userProto.ReqSignUp{
		Username: username,
		Password: passwd,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, module.ApiResp{
		ErrorNo:  int64(resp.Code),
		ErrorMsg: resp.Message,
	})
}

//// GetHostHandler : 获取服务器列表
//func GetHostHandler(c *gin.Context) {
//	appVersion := c.Request.FormValue("app_version")
//	appName := c.Request.FormValue("app_name")
//
//	resp, err := userCli.GetHost(context.TODO(), &userProto.ReqHost{
//		AppName:    appName,
//		AppVersion: appVersion,
//	})
//
//	if err != nil {
//		log.Println(err.Error())
//		c.Status(http.StatusInternalServerError)
//		return
//	}
//
//	c.JSON(http.StatusOK, module.ApiResp{
//		ErrorNo:  int64(resp.Code),
//		ErrorMsg: resp.Message,
//		Data:     resp.Host,
//	})
//}
