package model

import "xj_web_server/db"

type RecordExchange struct {
	ID              int     `json:"record_id"`
	UserID          int     `json:"-"`
	Amount          float64 `json:"amount"`
	PayCost         float64 `json:"-"`
	CurAmount       float64 `json:"-"`
	ChangeCoinRate  int     `json:"-"`
	ExchangeAccount string  `json:"exchange_account"`
	ExchangeWay     int     `json:"exchange_way"`
	ApplyStatus     int     `json:"apply_status"`
	ExchangeIP      string  `json:"-"`
	ApplyDate       string  `json:"apply_date"`
	HandlerDate     string  `json:"handler_date"`
	Remarks         string  `json:"-"`
	OperUserID      int     `json:"-"`
}

func (RecordExchange) TableName() string {
	return "RecordExchangeInfo"
}

func GetRecordExchangeByUid(uid, size, page int) ([]RecordExchange, error) {

	var data = make([]RecordExchange, 0)
	err := db.GetDBTreasure().Where(" UserID = ? ", uid).Limit(size, page*(size-1)).Find(&data)

	return data, err

}
