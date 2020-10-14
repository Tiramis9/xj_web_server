package model

import (
	"fmt"
	"log"
	redis "xj_web_server/cache"
	"xj_web_server/config"
	"xj_web_server/db"
	"testing"
)

func TestIsSensitiveWords(t *testing.T) {
	config.InitConfig("/../config/config.yml")
	err := redis.InitRedis(config.GetRedis())
	if err != nil {
		log.Fatalf("redis init err %v", err)
		return
	}
	initDB, err := db.InitDB(config.GetDb())
	if err != nil {
		log.Fatalf("db init err %v", err)
		return
	}
	defer func() {
		err := initDB.Close()
		if err != nil {
			log.Fatalf("db close err %v", err)
		}
	}()
	words, err := IsSensitiveWords("root")
	fmt.Println(words,err)
}