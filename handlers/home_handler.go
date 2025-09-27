package handlers

import (
	"fmt"
	"net/http"

	"gitee.com/wwgzr/blog/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 文章处理组
type HomeHanders struct {
	db *gorm.DB
}

func NewHomeHanders(db *gorm.DB) *HomeHanders {
	return &HomeHanders{db: db}
}

func (h *HomeHanders) ShowHome(c *gin.Context) {
	var articles []models.Article
	h.db.Find(&articles)
	var titles []string
	for _, article := range articles {
		titles = append(titles, article.Title)
		fmt.Println(article.Title)
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":  "文章列表",
		"titles": titles,
	})
}
