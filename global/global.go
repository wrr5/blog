package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	// 全局数据库变量
	DB  *gorm.DB
	RDB *redis.Client
)
