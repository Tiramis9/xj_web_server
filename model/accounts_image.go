package model

import "xj_web_server/db"

//CREATE TABLE `AccountsImage` (
//`UserID` int(11) NOT NULL COMMENT '用户标识',
//`RoleID` int(11) NOT NULL DEFAULT '1' COMMENT '角色标识（客户端数据，从1开始）',
//`SuitID` int(11) NOT NULL DEFAULT '1' COMMENT '套装标识（客户端数据，从1开始）',
//`PhotoFrameID` int(11) NOT NULL DEFAULT '1' COMMENT '头相框标识（客户端数据，从1开始）',
//PRIMARY KEY (`UserID`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='人物形象设置表';

type AccountsImage struct {
	UserID       int `json:"user_id"`        // 用户id
	RoleID       int `json:"role_id"`        // 角色标识
	SuitID       int `json:"suit_id"`        // 套装标识
	PhotoFrameID int `json:"photo_frame_id"` // 头像框标识
}

// 修改 角色 套装
func UpdateUserImageRoleByUid(uid, roleId, suitId int) error {

	accountsImage := new(AccountsImage)
	accountsImage.RoleID = roleId
	accountsImage.SuitID = suitId

	_, err := db.GetDB().Where(" UserID = ? ", uid).Update(accountsImage)

	return err
}

// 修改  头像框
func UpdateUserImageFrameByUid(uid, photoFrameID int) error {

	accountsImage := new(AccountsImage)
	accountsImage.PhotoFrameID = photoFrameID

	_, err := db.GetDB().Where(" UserID = ? ", uid).Update(accountsImage)

	return err
}
