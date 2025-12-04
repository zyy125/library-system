package database

import (
	"log"
	"time"
	"library-system/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL(dsn string) {
	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("MySQL 连接失败: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取底层 DB 失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
    sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
    sqlDB.SetConnMaxLifetime(time.Hour) // 连接最长可复用时间

	log.Println("MySQL 已连接")
}

func MigrateSQL() {
	err := DB.AutoMigrate(&model.Book{}, &model.BorrowRecord{}, &model.Category{}, &model.User{})
	if err != nil {
		log.Fatalf("MySQL自动迁移失败: %v", err)
	}
}
