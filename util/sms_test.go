package util

import "testing"

func TestSendSms(t *testing.T) {
	//SendSms(REGISTERED,"18822855252","123456")
	SendSms(ChangePassword,"18822855252","123456")
	//SendSms(BingPhoneNumber,"18822855252","123456")
}
