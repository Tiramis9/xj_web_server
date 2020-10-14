package db

import (
	"fmt"
	"xj_web_server/config"
	"testing"
)

type AccountsInfo struct {
	//`UserID` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户标识',
	//`FaceID` smallint(6) NOT NULL DEFAULT '0' COMMENT '头像标识',
	//`Accounts` varchar(32) NOT NULL COMMENT '用户帐号',
	UserID   int64
	FaceID   int64
	Accounts string
}

type User struct {
	Uid           int    `json:"uid" xorm:"UserID"`
	Mobile        string `json:"mobile" xorm:"RegisterMobile"`
	Token         string `json:"token" xorm:"-"`
	LoginPassword string `json:"-" xorm:"LogonPass"`
	Account string `json:"account" xorm:"Accounts"`
	CodeKey string `json:"-" xorm:"CodeKey"`
	FaceID  int `json:"face_id" xorm:"FaceID"`
	Nickname string `json:"nickname" xorm:"NickName"`
}

func (User) TableName() string {
	return "AccountsInfo"
}

func TestInitDB(t *testing.T) {
	config.InitConfig("/../config/config.yml")
	db, err := InitDB(config.GetDb())
	if err != nil {
		t.Fatalf("init db err %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("db close err %v", err)
		}
	}()
	//var accountsInfo []AccountsInfo
	//err = db.Where("UserID = ?", 30).Find(&accountsInfo)
	//if err != nil {
	//	t.Fatalf("db find err %v", err)
	//}
	//t.Logf("accountsInfo %v", accountsInfo)

	//user := new(User)
	var x []User
	//b, err := db.Where(" RegisterMobile = ?", "18822855256").Get(user)
	err = db.Limit(2,5,2).Find(&x)
	fmt.Println(err,x)

}
