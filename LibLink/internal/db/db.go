package db

import (
	"fmt"
	"liblink/internal/models/archive"
	"liblink/internal/models/system"
	"liblink/internal/models/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(host string, username string, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/liblink?charset=utf8&parseTime=True", username, password, host)
	return initDB(dsn)
}

func initDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别为 Info
	})
	if err != nil {
		fmt.Printf("%t\n", err)
		return nil, err
	}
	// 自动迁移
	err = db.AutoMigrate(
		&user.User{},
		&system.Notification{},
		&archive.Folder{},
		&archive.Archive{},
	)
	fmt.Printf("test db init\n")
	if err != nil {
		return nil, err
	}
	return db, err
}
