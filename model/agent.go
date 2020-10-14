package model

type RecordGame struct {
	UserID     int     `xorm:"UserID" json:"uid"`
	NickName   string  `xorm:"NickName" json:"nickname"`
	KindID     int     `xorm:"KindID" json:"kind_id"`
	KindName   string  `xorm:"KindName" json:"kind_name"`
	TotalScore float64 `xorm:"TotalScore" json:"total_score"`
}

type RecordGameAmount struct {
	UserId         int     `json:"uid"`
	ExchangeAmount float64 `json:"exchange_amount"`
	OrderAmount    float64 `json:"order_amount"`
	PresentAmount  float64 `json:"present_amount"`
}

type RecordDate struct {
	DateID     int     `xorm:"DateID" json:"date_id"`
	KindID     int     `xorm:"KindID" json:"kind_id"`
	KindName   string  `xorm:"KindName" json:"kind_name"`
	TotalScore float64 `xorm:"TotalScore" json:"total_score"`
}

type RecordGameDateAmount struct {
	DateID         int     `json:"date_id"`
	ExchangeAmount float64 `json:"exchange_amount"`
	OrderAmount    float64 `json:"order_amount"`
	PresentAmount  float64 `json:"present_amount"`
}

type RecordGameUserInfo struct {
	UserId      int     `json:"uid"`
	NickName    string  `json:"nickname"`
	Diamond     float64 `json:"diamond"`
	CollectDate string  `json:"collect_date"`
}

type RecordDiamondChangeLog struct {
	UserId        int     `json:"uid"`
	NickName      string  `json:"nickname"`
	CapitalTypeID int     `json:"capital_type_id"`
	CapitalAmount float64 `json:"capital_amount"`
	LastAmount    float64 `json:"last_amount"`
	LogDate       string  `json:"log_date"`
}
