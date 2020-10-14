package util

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

const (
	REGISTERED      = "SMS_193511188" //注册
	ChangePassword  = "SMS_193511189" //修改密码
	BingPhoneNumber = "SMS_193512698" //绑定手机号码
)

func SendSms(modeCode, number, code string) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI4GEsWQyoozM2YAm6K6Si", "XwhOfPBpGFb8g5TcC6uhSk2ry1WTIO")

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = number
	request.SignName = "XJGame"
	request.TemplateCode = modeCode
	request.TemplateParam = `{"code":` + code + `}`

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}
