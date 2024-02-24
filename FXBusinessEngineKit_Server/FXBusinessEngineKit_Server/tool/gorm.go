package tool

import (
	"FXBusinessEngineKit_Server/configuration"
	"FXBusinessEngineKit_Server/log"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var FXGorm *gorm.DB

func InitConnectDatabase() {
	dsn := configuration.FXDatabaseServer
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.RecordLog(log.Err, fmt.Sprintf("[数据库]数据库连接失败%s", err))
		return
	}
	FXGorm = db
	log.RecordLog(log.Msg, "数据库连接成功")
}
