package database

import (
	"database/sql"
	"entry_task/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(mysql.Open(config.ProjectConfig.Db.Dsn), &gorm.Config{Logger: newGormLog(nil, nil)})
	if nil != err {
		panic("fail to connect database,connect failed, err:" + err.Error())
	}
	sqldb, err := db.DB()
	configConnectPool(sqldb)
}

func GetDB(requestId *string) *gorm.DB {
	if requestId == nil {
		return db
	} else {
		return db.Session(&gorm.Session{Logger: newGormLog(requestId, nil)})
	}
}

func configConnectPool(sqldb *sql.DB) {
	sqldb.SetMaxIdleConns(config.ProjectConfig.Db.MaxIdle)
	sqldb.SetConnMaxIdleTime(time.Duration(config.ProjectConfig.Db.MaxIdleTime) * time.Second)
	sqldb.SetMaxOpenConns(config.ProjectConfig.Db.MaxOpen)
}
