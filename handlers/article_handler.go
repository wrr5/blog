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

func (h *ArticleHanders) ShowArticleList(c *gin.Context) {
	var articles []models.Article
	h.db.Find(&articles)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":    "文章列表",
		"articles": articles,
	})
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

	c.Redirect(http.StatusFound, "/articles")
}

func (h *ArticleHanders) ShowArticleDetail(c *gin.Context) {
	var article models.Article
	id := c.Param("id")
	h.db.First(&article, id)

	c.HTML(http.StatusOK, "article.html", gin.H{
		"article": article,
	})
}

func (h *ArticleHanders) ShowArticleEdit(c *gin.Context) {
	id := c.Param("id")
	var article models.Article
	h.db.First(&article, id)
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"article": article,
	})
}

func (h *ArticleHanders) UpdateArticle(c *gin.Context) {
	type UpdateArticleRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	id := c.Param("id")
	var article models.Article
	var req UpdateArticleRequest

	// 获取文章与修改后的内容
	h.db.First(&article, id)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 修改文章
	article.Title = req.Title
	article.Content = req.Content
	h.db.Save(&article)

	c.JSON(http.StatusOK, gin.H{
		"message":  "修改成功",
		"redirect": "/articles/" + id, // 跳转URL
		"article":  article,           // 返回更新后的文章数据
	})
}

func (h *ArticleHanders) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article
	h.db.Delete(&article, id)
	c.JSON(http.StatusOK, gin.H{
		"message":  "删除成功",
		"redirect": "/articles/" + id, // 跳转URL
		"article":  article,           // 返回更新后的文章数据
	})
}
