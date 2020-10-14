package model

import "xj_web_server/db"

type Wheel struct {
	ItemIndex int `json:"item_index"`
	ItemQuota float64 `json:"item_quota"`
}

type WheelRecord struct {
	RecordID int `json:"record_id"`
	UserID int `json:"-"`
	ItemIndex int `json:"item_index"`
	ScoreType int `json:"score_type"`
	ItemQuota int `json:"item_quota"`
	CollectDate string `json:"collect_date"`
}

func (WheelRecord) TableName() string {
	return "RecordTurntable"
}

func GetWheelRecord(uid , wheelType, size , page  int) ([]WheelRecord,error) {
	var wheelRecord []WheelRecord

	err := db.GetDBRecord().Where(" UserID = ? AND ScoreType = ? ", uid, wheelType).Desc("CollectDate").Limit(size,page * (size - 1)).Find(&wheelRecord)

	return wheelRecord,err
}