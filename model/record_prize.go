package model

import "xj_web_server/db"

type RecordPrizeInfo struct {
	ID         int     `json:"id"`
	NickName   string  `json:"nick_name"`
	GameName   string  `json:"game_name"`
	RoomName   string  `json:"room_name"`
	Describe   string  `json:"describe"`
	CellScore  float64 `json:"cell_score"`
	Multiply   float64 `json:"multiply"`
	Score      float64 `json:"score"`
	RecordDate string  `json:"record_date"`
}

func (RecordPrizeInfo) TableName() string {
	return "RecordPrizeInfo"
}

func GetRecordPrizeInfo(size, page int) ([]RecordPrizeInfo, error) {

	var data = make([]RecordPrizeInfo, 0)
	err := db.GetDBTreasure().Limit(size, size*(page-1)).Find(&data)

	return data, err

}
