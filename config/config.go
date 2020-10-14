package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	yaml2 "gopkg.in/yaml.v2"
)

// DB 数据库配置
type Db struct {
	EnableLog          bool   `yaml:"enable_log" json:"enable_log"`
	Dialect            string `yaml:"dialect" json:"dialect"`
	Host               string `yaml:"host" json:"host"`
	User               string `yaml:"user" json:"user"`
	PassWd             string `yaml:"pass" json:"pass"`
	Db                 string `yaml:"db" json:"db"`
	MaxOpenConnections int    `yaml:"max_open_connections" json:"max_open_connections"`
	MaxIdleConnections int    `yaml:"max_idle_connections" json:"max_idle_connections"`
}

// Service 服务端配置
type Service struct {
	Mode        string `yaml:"mode"`
	Port        string `yaml:"port"`
	TCPPort     string `yaml:"tcp_port"`
	ServiceUrl  string `yaml:"service_url"`
	RpcUrl      string `yaml:"rpc_url"`
	MaxPageSize int    `yaml:"max_page_size"`
	// 超时时间
	AppReadTimeout  time.Duration `yaml:"app_read_timeout"`
	AppWriteTimeout time.Duration `yaml:"app_write_timeout"`
}

// redis
type Redis struct {
	Host   string `yaml:"host"`
	PassWd string `yaml:"pass"`
	Db     int    `yaml:"db"`
}

type LogConfig struct {
	Path string `yaml:"path"`
}

type WssConfig struct {
	HeartbeatTime int64 `yaml:"heart_beat_time"`
}

// Config 配置
type Config struct {
	Service    Service   `yaml:"service"`
	DB         Db        `yaml:"db"`
	DBTreasure Db        `yaml:"db_treasure"`
	DBPlatform Db        `yaml:"db_platform"`
	DBRecord   Db        `yaml:"db_record"`
	Redis      Redis     `yaml:"redis"`
	Log        LogConfig `yaml:"log"`
	Wss        WssConfig `yaml:"wss"`
}

func GetDb() Db {
	return config.DB
}

func GetDBTreasure() Db {
	return config.DBTreasure
}

func GetDBPlatform() Db {
	return config.DBPlatform
}

func GetDBRecord() Db {
	return config.DBRecord
}

func GetService() Service {
	return config.Service
}

func GetRedis() Redis {
	return config.Redis
}

func GetWss() WssConfig {
	return config.Wss
}

// GetLog 日志配置
func GetLog() LogConfig {
	return config.Log
}

/**
 * @description:
 * @param {type}
 * @return:
 */
var config Config

// InitConfig 初始化config
func InitConfig(path string) {
	pathStr, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	yamlFile, err := ioutil.ReadFile(pathStr + path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml2.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
	}
}
