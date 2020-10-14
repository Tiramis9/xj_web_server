package model

type Record struct {
	RecordID  int     `json:"record_id"`
	UserID    int     `json:"uid"`
	UserScore float64 `json:"user_score"`
	KindName  string  `json:"kind_name"`
	GameName  string  `json:"game_name"`
	EnterTime string  `json:"enter_time"`
	LeaveTime string  `json:"leave_time"`
}

func (t *Record) TableName() string {

	return "record"
}
