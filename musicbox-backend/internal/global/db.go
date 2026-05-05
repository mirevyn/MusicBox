package global

import (
	"fmt"
	"log"
	"musicbox-backend/internal/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 从配置中获取 MySQL 连接串
	dsn := config.Conf.MySQL.DSN

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("获取数据库连接池失败: ", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("数据库连通性检查失败: ", err)
	}

	fmt.Println("MySQL 连接成功!")
}
