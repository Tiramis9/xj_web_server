package model

import "xj_web_server/db"

type AgentRoyaltyTakeRecord struct {
	Id            int64   `json:"-"`
	UserID        int64   `json:"-"`
	RoyaltyAmount float64 `json:"royalty_amount"`
	TakeDate      string  `json:"take_date"`
	ClientIP      string  `json:"-"`
}

func (AgentRoyaltyTakeRecord) TableName() string {
	return "AgentRoyaltyTakeRecord"
}

func GetAgentRoyaltyTakeRecordByUid(uid int) ([]AgentRoyaltyTakeRecord, error) {

	var royalty = make([]AgentRoyaltyTakeRecord,0)

	err := db.GetDBTreasure().Where(" UserID = ?", uid).Find(&royalty)

	return royalty, err
}
