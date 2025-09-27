package handlers

import (
	"net/http"

	"gitee.com/wwgzr/blog/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 文章处理组
type ArticleHanders struct {
	db *gorm.DB
}

func NewArticleHanders(db *gorm.DB) *ArticleHanders {
	return &ArticleHanders{db: db}
}

func (h *ArticleHanders) ShowCreateArticlePage(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", gin.H{})
}

func (h *ArticleHanders) CreateArticle(c *gin.Context) {
	// 获取表单数据
	title := c.PostForm("title")
	content := c.PostForm("content")

	// 创建文件对象
	article := models.Article{
		Title:   title,
		Content: content,
	}

	result := h.db.Create(&article)
	if result.Error != nil {
		c.HTML(http.StatusOK, "new.html", gin.H{
			"error": result.Error,
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}
