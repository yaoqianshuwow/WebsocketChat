package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"WebsocketChat/internal/config"
	"WebsocketChat/internal/model"
	"WebsocketChat/pkg/zlog"
)

var GormDB *gorm.DB

func init() {
	// 从配置文件获取数据库连接信息
	cfg := config.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
		cfg.MysqlConfig.User, cfg.MysqlConfig.Password, cfg.MysqlConfig.Host, cfg.MysqlConfig.Port, cfg.MysqlConfig.DatabaseName)
	var err error
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zlog.Fatal(err.Error())
	}
	err = GormDB.AutoMigrate(&model.UserInfo{}, &model.GroupInfo{}, &model.UserContact{}, &model.Session{}, &model.ContactApply{}, &model.Message{}) // 自动迁移，如果没有建表，会自动创建对应的表
	if err != nil {
		zlog.Fatal(err.Error())
	}
}
