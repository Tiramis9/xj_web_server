/*
 * @Author: yhlyl
 * @Date: 2019-11-05 16:37:38
 * @LastEditTime: 2019-11-06 11:28:03
 * @LastEditors: yhlyl
 * @Description:
 * @FilePath: /xj_web_server/cache/redis_test.go
 * @https://github.com/android-coco
 */
package redis

import (
	"fmt"
	"xj_web_server/config"
	"testing"
)

func TestInitRedis(t *testing.T) {
	var redisConfig = config.Redis{Host: "127.0.0.1:6379", PassWd: "uJREJW9DNIk2H3I96ayz", Db: 0}
	InitRedis(redisConfig)
	err := GetRedisDb().LPush("list", "fsad111").Err()
	lLen, _ := GetRedisDb().LLen("list").Result()
	s, err := GetRedisDb().LRange("list", 0, lLen).Result()

	x := GetRedisDb().Exists("19999").Val()
	fmt.Print(x)
	get,err := GetRedisDb().Get("user:login:token:31:").Result()
	fmt.Println(s,get, err)
}
