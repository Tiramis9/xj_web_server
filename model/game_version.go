package model

import (
	"github.com/go-xorm/xorm"
)

type GameVersion struct {
	KindID    int64  `json:"kind_id"`
	Version   int64  `json:"version"`
	PokerName string `json:"poker_name"`
	Hash      string `json:"hash"`
	Size      int64  `json:"size"`
	Platform  string `json:"platform"`
}

func (GameVersion) TebleName() string {
	return "GameVersion"
}

type SiteConfiginfo struct {
	ConfigKey string `json:"config_key"`
	Field1    string `json:"field_1"`
	Field2    string `json:"field_2"`
	Field3    string `json:"field_3"`
	Field4    string `json:"field_4"`
}

func (SiteConfiginfo) TebleName() string {
	return "SiteConfiginfo"
}

func GetVersionConf(db *xorm.Engine, configKey string) ([]SiteConfiginfo, error) {
	var gameVersions []SiteConfiginfo
	err := db.Where(" ConfigKey = ? ", configKey).Find(&gameVersions)

	if err != nil {
		return nil, err
	}

	return gameVersions, nil
}

func GetAllGameVersion(db *xorm.Engine, platform string) ([]GameVersion, error) {
	var gameVersions []GameVersion

	err := db.Where(" Platform = ? ", platform).Find(&gameVersions)

	if err != nil {
		return nil, err
	}

	return gameVersions, nil
}

func SetGameVersion(db *xorm.Engine, size int64, name, hash, platform string) (*GameVersion, error) {
	gameVersion := new(GameVersion)

	_, err := db.Where(" PokerName = ? AND Platform = ? ", name, platform).Get(gameVersion)

	if err != nil {
		return nil, err
	}

	newGameVersion := new(GameVersion)

	newGameVersion.PokerName = name
	newGameVersion.Hash = hash
	newGameVersion.Size = size
	newGameVersion.Platform = platform

	if gameVersion.KindID == 0 {
		newGameVersion.Version = 1
		_, err = db.Insert(newGameVersion)
	} else {
		newGameVersion.Version = gameVersion.Version + 1
		_, err = db.Where(" KindID = ? ", gameVersion.KindID).Update(newGameVersion)
	}

	return newGameVersion, err

}
