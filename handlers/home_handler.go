package handlers

import (
	// "net/http"

	// "gitee.com/wwgzr/blog/models"
	// "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 文章处理组
type HomeHanders struct {
	db *gorm.DB
}

func NewHomeHanders(db *gorm.DB) *HomeHanders {
	return &HomeHanders{db: db}
}
