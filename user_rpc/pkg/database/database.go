package database

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"user_rpc/pkg/config"
	"user_rpc/pkg/logger"
)

var DB *gorm.DB
var SqlDB *sql.DB

// SetupDatabase 初始化数据库
func SetupDatabase() {
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		config.GetString("DB_USERNAME"),
		config.GetString("DB_PASSWORD"),
		config.GetInt("DB_HOST"),
		config.Get("DB_PORT"),
		config.Get("DB_DATABASE"),
		config.Get("DB_CHARSET"),
	)
	dbConfig := mysql.New(mysql.Config{DSN: dsn})
	DB, err = gorm.Open(dbConfig)
	if err != nil {
		logger.Error("database", "open db dsn error", err)
		panic(err)
	}

	SqlDB, err = DB.DB()
	if err != nil {
		logger.Error("database", "get sqldb error", err)
		panic(err)
	}

	SqlDB.SetMaxOpenConns(config.GetInt("DB_MAX_OPEN_CONN"))
	SqlDB.SetMaxIdleConns(config.GetInt("DB_MAX_IDLE_CONN"))
	SqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("DB_MAX_LIFETIME")) * time.Hour)

	logger.Debug("database", "conn success")
}

// Close 关闭数据库连接
func Close() {
	if SqlDB == nil {
		return
	}
	if err := SqlDB.Close(); err != nil {
		logger.Error("database", "close db conn err", err)
		panic(err)
	}
	logger.Debug("database", "close success")
}
