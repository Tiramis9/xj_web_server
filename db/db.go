package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"xj_web_server/config"
	"xj_web_server/util"
	"time"
	"xorm.io/core"
)

var db *xorm.Engine
var dbTreasure *xorm.Engine
var dBPlatform *xorm.Engine
var dBRecord *xorm.Engine

type dbLogger struct {
}

func (logger dbLogger) Level() core.LogLevel {
	return core.LOG_INFO
}

func (logger dbLogger) SetLevel(l core.LogLevel) {
}

func (logger dbLogger) ShowSQL(show ...bool) {
}

func (logger dbLogger) IsShowSQL() bool {
	return true
}

func (logger dbLogger) Debug(v ...interface{}) {
	util.Logger.Debug(v)
}

func (logger dbLogger) Debugf(format string, v ...interface{}) {
	util.Logger.Debugf(format, v)
}

func (logger dbLogger) Error(v ...interface{}) {
	util.Logger.Error(v)
}

func (logger dbLogger) Errorf(format string, v ...interface{}) {
	util.Logger.Errorf(format, v)
}

func (logger dbLogger) Info(v ...interface{}) {
	util.Logger.Info(v)
}

func (logger dbLogger) Warn(v ...interface{}) {
	util.Logger.Warn(v)
}

func (logger dbLogger) Warnf(format string, v ...interface{}) {
	util.Logger.Warnf(format, v)
}

func (logger dbLogger) Infof(format string, v ...interface{}) {
	util.Logger.Infof(format, v)
}

// Connect 初始化 DB连接
func InitDB(configDb config.Db) (*xorm.Engine, error) {
	args := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configDb.User, configDb.PassWd, configDb.Host, configDb.Db)
	var err error
	db, err = xorm.NewEngine(configDb.Dialect, args)
	if err != nil {
		log.Fatalf("init db err %v \n", err)
	}
	db.ShowSQL(configDb.EnableLog)
	db.SetLogger(dbLogger{})
	db.SetLogLevel(core.LOG_INFO)
	// 结构体名称和对应的表名称以及结构体field名称与对应的表字段名称相同的命名
	db.SetMapper(core.SameMapper{})
	//用于设置最大打开的连接数，默认值为0表示不限制。
	db.SetMaxOpenConns(configDb.MaxOpenConnections)
	//设置连接池的空闲数大小
	db.SetMaxIdleConns(configDb.MaxIdleConnections)
	syncTable()
	return db, nil
}

// Connect 初始化 DB连接
func InitDBTreasure(configDb config.Db) (*xorm.Engine, error) {
	args := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configDb.User, configDb.PassWd, configDb.Host, configDb.Db)
	var err error
	dbTreasure, err = xorm.NewEngine(configDb.Dialect, args)
	if err != nil {
		log.Fatalf("init db err %v \n", err)
	}
	dbTreasure.ShowSQL(configDb.EnableLog)
	dbTreasure.SetLogger(dbLogger{})
	dbTreasure.SetLogLevel(core.LOG_INFO)
	// 结构体名称和对应的表名称以及结构体field名称与对应的表字段名称相同的命名
	dbTreasure.SetMapper(core.SameMapper{})
	//用于设置最大打开的连接数，默认值为0表示不限制。
	dbTreasure.SetMaxOpenConns(configDb.MaxOpenConnections)
	//设置连接池的空闲数大小
	dbTreasure.SetMaxIdleConns(configDb.MaxIdleConnections)
	syncTable()
	return dbTreasure, nil
}

// Connect 初始化 DB连接
func InitDBPlatform(configDb config.Db) (*xorm.Engine, error) {
	args := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configDb.User, configDb.PassWd, configDb.Host, configDb.Db)
	var err error
	dBPlatform, err = xorm.NewEngine(configDb.Dialect, args)
	if err != nil {
		log.Fatalf("init db err %v \n", err)
	}
	dBPlatform.ShowSQL(configDb.EnableLog)
	dBPlatform.SetLogger(dbLogger{})
	dBPlatform.SetLogLevel(core.LOG_INFO)
	// 结构体名称和对应的表名称以及结构体field名称与对应的表字段名称相同的命名
	dBPlatform.SetMapper(core.SameMapper{})
	//用于设置最大打开的连接数，默认值为0表示不限制。
	dBPlatform.SetMaxOpenConns(configDb.MaxOpenConnections)
	//设置连接池的空闲数大小
	dBPlatform.SetMaxIdleConns(configDb.MaxIdleConnections)
	dBPlatform.DatabaseTZ = time.Local // 必须
	dBPlatform.TZLocation = time.Local // 必须
	syncTable()
	return dBPlatform, nil
}

// Connect 初始化 DB连接
func InitDBRecord(configDb config.Db) (*xorm.Engine, error) {
	args := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configDb.User, configDb.PassWd, configDb.Host, configDb.Db)
	var err error
	dBRecord, err = xorm.NewEngine(configDb.Dialect, args)
	if err != nil {
		log.Fatalf("init db err %v \n", err)
	}
	dBRecord.ShowSQL(configDb.EnableLog)
	dBRecord.SetLogger(dbLogger{})
	dBRecord.SetLogLevel(core.LOG_INFO)
	// 结构体名称和对应的表名称以及结构体field名称与对应的表字段名称相同的命名
	dBRecord.SetMapper(core.SameMapper{})
	//用于设置最大打开的连接数，默认值为0表示不限制。
	dBRecord.SetMaxOpenConns(configDb.MaxOpenConnections)
	//设置连接池的空闲数大小
	dBRecord.SetMaxIdleConns(configDb.MaxIdleConnections)
	dBPlatform.DatabaseTZ = time.Local // 必须
	dBPlatform.TZLocation = time.Local // 必须
	syncTable()
	return dBRecord, nil
}

func GetDB() *xorm.Engine {
	return db
}

func GetDBTreasure() *xorm.Engine {
	return dbTreasure
}

func GetDBPlatform() *xorm.Engine {
	return dBPlatform
}

func GetDBRecord() *xorm.Engine {
	return dBRecord
}

func syncTable() {
	//err := db.Sync2(&model.User{}, &model.Community{}, &model.Contact{})
	//if err != nil {
	//	log.Fatalf("init db sync table err %v \n", err)
	//}
}
