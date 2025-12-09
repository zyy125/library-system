package database

import (
	"log"
	"time"
	"fmt"
	"library-system/model"
	"library-system/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MigrateSQL(db *gorm.DB) error{
	err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Book{},
		&model.BorrowRecord{},
	)
	if err != nil {
		return fmt.Errorf("MySQL自动迁移失败: %v", err)
	}
	return nil
}

func InitMySQL() (*gorm.DB, error){
	sqlCfg := config.Load()

	db, err := gorm.Open(mysql.Open(sqlCfg.DSN), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("MySQL 连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层 DB 失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
    sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
    sqlDB.SetConnMaxLifetime(time.Hour) // 连接最长可复用时间

	log.Println("MySQL 已连接")

	if err = MigrateSQL(db); err != nil {
		return nil, fmt.Errorf("获取底层 DB 失败: %v", err)
	}
	
	return db, nil
}

