package model

import "xj_web_server/db"

type AgentRoyaltyLevel struct {
	LevelID       int64   `json:"level_id"`
	LevelNum      int64   `json:"level_num"`
	MinTotalMoney float64 `json:"min_total_money"`
	PercentValue  float64 `json:"percent_value"`
}

func (AgentRoyaltyLevel) TableName() string {
	return "AgentRoyaltyLevel"
}

func GetAllAgentRoyaltyLevel() ([]AgentRoyaltyLevel, error) {

	var royalty []AgentRoyaltyLevel

	err := db.GetDBTreasure().Find(&royalty)

	return royalty, err
}

type SystemStatusInfo struct {
	StatusValue int64 `json:"status_value"`
}

func (SystemStatusInfo) TableName() string {
	return "SystemStatusInfo"
}

func GetSystemStatusInfo() ([]SystemStatusInfo, error) {

	var royalty []SystemStatusInfo

	err := db.GetDBPlatform().Where("StatusName = ? OR StatusName = ? OR StatusName = ?", "RoyaltyRechargeAmount", "RoyaltyTeamCount", "RechargeReward").Find(&royalty)

	return royalty, err
}
