package model

import (
	"xj_web_server/db"
)

type News struct {
	NewsID         int        `json:"news_id"`
	Subject        string     `json:"subject"`
	IsEnable       int        `json:"is_enable"`
	IsDelete       int        `json:"is_delete"`
	Body           string     `json:"body"`
	FormattedBody  string     `json:"formatted_body"`
	CreateDate     LocalTime `json:"issue_date"`
	LastModifyDate LocalTime `json:"last_modify_date"`
	SortID         int        `json:"sort_id"`
}

func (News) TableName() string {
	return "NewsInfo"
}

func GetNewsInfo() ([]News, error) {

	var news = make([]News, 0)

	err := db.GetDBPlatform().Where(" IsEnable = ? AND IsDelete = ? ", 1, 0).OrderBy("SortID").Find(&news)

	return news, err
}
