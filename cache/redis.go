package redis

import (
	"github.com/go-redis/redis/v7"
	"xj_web_server/config"
	"time"
)

var db *redis.Client

func InitRedis(redisConfig config.Redis) error {
	db = redis.NewClient(&redis.Options{
		Addr:         redisConfig.Host,
		Password:     redisConfig.PassWd,
		DB:           redisConfig.Db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	_, err := db.Ping().Result()
	return err
}

// GetRedisDb 获取Redis连接实例
func GetRedisDb() *redis.Client {
	return db
}
