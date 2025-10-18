package tools

import (
	"fmt"
	"log"

	"gitee.com/wwgzr/blog/config"
	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password)
	db, error := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if error != nil {
		log.Fatal("数据库连接失败：", error)
	}

	// 自动迁移（如果表不存在则创建, 已存在则检查有无新增字段，不会修改字段名和删除字段）
	err := db.AutoMigrate(&models.Article{}, &models.User{}, &models.Category{}, &models.Tag{}, &models.Comment{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
	log.Println("数据库连接并迁移成功!")

	global.DB = db
	return db
}
