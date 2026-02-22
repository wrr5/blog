package tools

import (
	"fmt"
	"log"
	"os"

	"gitee.com/wwgzr/blog/config"
	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/blog?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port)
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

	var rdb *redis.Client
	redis_cfg := config.AppConfig.Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redis_cfg.Host, redis_cfg.Port), // 例如 "localhost:6379"
		Password: os.Getenv("REDIS_PASSWORD"),                          // 没有密码则留空
		DB:       0,                                                    // 使用默认DB
		PoolSize: 10,                                                   // 连接池大小
	})
	global.RDB = rdb

	return db
}
