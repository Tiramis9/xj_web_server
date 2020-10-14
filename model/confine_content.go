package model

import (
	"xj_web_server/db"
	"time"
)

type ConfineContent struct {
	ConfineString  string     //限制字符
	EnjoinOverDate *time.Time //过期时间
	CollectDate    *time.Time //收集日期
	CollectNote    string     //备注
}

func (ConfineContent) TableName() string {
	return "ConfineContent"
}

func IsSensitiveWords(word string) (bool, error) {
	confineContent := new(ConfineContent)
	b, err := db.GetDB().Where(" ConfineString = ? ", word).Get(confineContent)
	if err != nil {
		return true, err
	}
	if b {
		// 永久有效
		if confineContent.EnjoinOverDate == nil && confineContent.ConfineString != "" {
			return true, nil
		}
		b, err = db.GetDB().Where(" ConfineString = ? AND EnjoinOverDate > NOW() ", word).Get(confineContent)
	}
	return b, err
}
