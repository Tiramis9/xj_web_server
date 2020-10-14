package model

import (
	"errors"
	"xj_web_server/db"
	"xj_web_server/httpserver/wss/proto"
	"strconv"
	"time"
)

type NewsInfo struct {
	NewsID        int32
	Subject       string
	Body          string
	FormattedBody string
	ClassID       int32
}


func GetUerByUid(uid int32) (*proto.Hall_S_Msg, error) {
	//SELECT info.*,game.GoldCoin,game.Diamond FROM xjaccountsdb.AccountsInfo info INNER JOIN xjtreasuredb.GameScoreInfo game ON info.UserID = game.UserID WHERE info.UserID=18;
	var announcement []NewsInfo
	db.GetDBPlatform().Find(&announcement)
	res, err := db.GetDB().QueryString("SELECT info.UserID,info.NickName,info.Gender,info.LevelNum,info.RegisterMobile,game.GoldCoin,game.Diamond,info.FaceID,image.RoleID,image.SuitID,image.PhotoFrameID,IFNULL(exchange.AccountOrCard,'') AS BinderCardNo FROM xjaccountsdb.AccountsInfo info INNER JOIN xjtreasuredb.GameScoreInfo game ON info.UserID = game.UserID INNER JOIN xjaccountsdb.AccountsImage image ON info.UserID = image.UserID LEFT JOIN xjaccountsdb.ExchangeAccount exchange ON info.UserID = exchange.UserID WHERE info.UserID=" + strconv.Itoa(int(uid)))
	if err != nil {
		return nil, err
	}
	if len(res) <= 0 {
		return nil, errors.New("用户不存在")
	}
	uidInt, _ := strconv.Atoi(res[0]["UserID"])
	levelNumInt, _ := strconv.Atoi(res[0]["LevelNum"])
	faceIDInt, _ := strconv.Atoi(res[0]["FaceID"])
	roleIDInt, _ := strconv.Atoi(res[0]["RoleID"])
	suitIDInt, _ := strconv.Atoi(res[0]["SuitID"])
	photoFrameIDInt, _ := strconv.Atoi(res[0]["PhotoFrameID"])
	genderInt, _ := strconv.Atoi(res[0]["Gender"])
	userGold, _ := strconv.ParseFloat(res[0]["GoldCoin"], 64)
	userDiamonds, _ := strconv.ParseFloat(res[0]["Diamond"], 64)
	var ans = make([]*proto.Announcement, 0)
	for _, v := range announcement {
		temp := &proto.Announcement{
			NewsID:        v.NewsID,
			Subject:       v.Subject,
			Body:          v.Body,
			FormattedBody: v.FormattedBody,
			ClassID:       v.ClassID,
		}
		ans = append(ans, temp)
	}
	return &proto.Hall_S_Msg{
		UserID:           int32(uidInt),
		NikeName:         res[0]["NikeName"],
		UserGold:         float32(userGold),
		UserDiamonds:     float32(userDiamonds),
		MemberOrder:      int32(levelNumInt),
		PhoneNumber:      res[0]["PhoneNumber"],
		BinderCardNo:     res[0]["BinderCardNo"],
		FaceID:           int32(faceIDInt),
		RoleID:           int32(roleIDInt),
		SuitID:           int32(suitIDInt),
		PhotoFrameID:     int32(photoFrameIDInt),
		Gender:           int32(genderInt),
		TimeStamp:        time.Now().Unix(),
		AnnouncementList: ans,
	}, nil
}
