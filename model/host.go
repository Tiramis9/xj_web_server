package model

import (
	"github.com/go-xorm/xorm"
)

type QpHost struct {
	ServerAddr   string `json:"server_addr"`
	WsServerAddr string `json:"ws_server_addr"`
}

func (QpHost) TableName() string {
	return "host"
}

func GetHost(db *xorm.Engine) ([]QpHost, error) {
	var hosts []QpHost

	err := db.Find(&hosts)

	if err != nil {
		return nil, err
	}

	return hosts, nil
}
