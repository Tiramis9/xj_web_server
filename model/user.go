package model

import (
	"xj_web_server/db"
	"xj_web_server/util"
)

//CREATE TABLE `user` (
//`id` bigint(64) NOT NULL AUTO_INCREMENT,
//`mobile` varchar(20) DEFAULT NULL,
//`passwd` varchar(40) DEFAULT NULL,
//`avatar` varchar(150) DEFAULT NULL,
//`sex` varchar(2) DEFAULT NULL,
//`nickname` varchar(20) DEFAULT NULL,
//`salt` varchar(10) DEFAULT NULL,
//`online` int(10) DEFAULT NULL,
//`token` varchar(40) DEFAULT NULL,
//`memo` varchar(140) DEFAULT NULL,
//`createat` datetime DEFAULT NULL,
//`uid` varbinary(20) DEFAULT NULL,
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
type User struct {
	Uid           int    `json:"uid" xorm:"UserID"`
	Mobile        string `json:"mobile" xorm:"RegisterMobile"`
	Token         string `json:"token" xorm:"-"`
	LoginPassword string `json:"-" xorm:"LogonPass"`
	Account       string `json:"account" xorm:"Accounts"`
	CodeKey       string `json:"-" xorm:"CodeKey"`
	FaceID        int    `json:"face_id" xorm:"FaceID"`
	Nickname      string `json:"nickname" xorm:"NickName"`
}

func (User) TableName() string {
	return "AccountsInfo"
}

func GetUserById(uid string) (User, error) {

	var user []User
	err := db.GetDB().Where(" uid = ? ", uid).Find(&user)

	return user[0], err
}

func GetUserByAccount(account string) (*User, bool, error) {
	user := new(User)

	b, err := db.GetDB().Where(" RegisterMobile = ? ", account).Get(user)

	return user, b, err
}

func SetUserPasswordByUid(uid int, password string) error {

	user := new(User)
	user.LoginPassword = password

	_, err := db.GetDB().Where(" UserID = ? ", uid).Update(user)

	return err
}

func UpdateUserFaceIdByUid(uid, faceId int) error {

	user := new(User)
	user.FaceID = faceId

	_, err := db.GetDB().Where(" UserID = ? ", uid).Update(user)

	return err
}

func UpdateUserNicknameByUid(uid int, nickname string) error {

	user := new(User)
	user.Nickname = nickname

	_, err := db.GetDB().Where(" UserID = ? ", uid).Update(user)

	return err
}

func IsExistUserInfo(uid int) bool {

	rows, err := db.GetDB().DB().DB.Query("CALL COM_SP_ISExistUserInfo(?,?)", uid, "1")
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				util.Logger.Errorf("rows 关闭错误 出错 err: %s ", err.Error())
			}
		}
	}()
	if err != nil {
		util.Logger.Errorf("IsExistUserInfo  查询存储过程 err: %s ", err.Error())
		return true
	}

	var errorCode int64
	var errorMsg string
	rows.Next()
	err = rows.Scan(&errorCode, &errorMsg)
	if err != nil {
		util.Logger.Errorf("IsExistUserInfo  查询存储过程 err: %s ", err.Error())
		return true
	}
	if errorCode != util.SuccessCode {
		return false
	}
	return true

}
