package handlers

import (
	"net/http"

	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"github.com/gin-gonic/gin"
)

// 文章处理组
type ArticleHanders struct {
}

func NewArticleHanders() *ArticleHanders {
	return &ArticleHanders{}
}

func (h *ArticleHanders) ShowArticleList(c *gin.Context) {
	var articles []models.Article
	global.DB.Find(&articles)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"title":    "文章列表",
		"articles": articles,
	})
}

func (h *ArticleHanders) ShowCreateArticlePage(c *gin.Context) {
	if username, ok := c.Get("username"); ok {
		c.HTML(http.StatusOK, "create.html", gin.H{
			"username": username,
		})
	}
}

func (h *ArticleHanders) CreateArticle(c *gin.Context) {
	// 获取表单数据
	title := c.PostForm("title")
	content := c.PostForm("content")

	// 获取作者对象
	var user models.User
	if username, ok := c.Get("username"); ok {
		result := global.DB.Where("username = ?", username).First(&user)
		if result.Error != nil {
			c.HTML(http.StatusOK, "create.html", gin.H{
				"error": result.Error,
			})
			return
		}

		// 创建文章对象
		article := models.Article{
			Title:   title,
			Content: content,
			User:    user,
		}

		result = global.DB.Create(&article)
		if result.Error != nil {
			c.HTML(http.StatusOK, "create.html", gin.H{
				"error": result.Error,
			})
			return
		}

		c.Redirect(http.StatusFound, "/articles")
	}
}

func (h *ArticleHanders) ShowArticleDetail(c *gin.Context) {
	var article models.Article
	articleID := c.Param("id")
	global.DB.Preload("User").First(&article, articleID)

	c.HTML(http.StatusOK, "article.html", gin.H{
		"user_id": c.GetUint("user_id"),
		"article": article,
	})
}

func (h *ArticleHanders) ShowArticleEdit(c *gin.Context) {
	// 获取文章id
	id := c.Param("id")
	var article models.Article
	global.DB.First(&article, id)

	// 获取当前用户登陆的id
	userID := c.GetUint("user_id")
	if article.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "无权操作此文章"})
		c.Abort()
		return
	}

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
	global.DB.First(&article, id)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 修改文章
	article.Title = req.Title
	article.Content = req.Content
	global.DB.Save(&article)

	c.JSON(http.StatusOK, gin.H{
		"message":  "修改成功",
		"redirect": "/articles/" + id, // 跳转URL
		"article":  article,           // 返回更新后的文章数据
	})
}

func (h *ArticleHanders) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article
	global.DB.First(&article, id)

	// 获取当前用户登陆的id
	userID := c.GetUint("user_id")
	if article.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "无权操作此文章"})
		c.Abort()
		return
	}

	global.DB.Delete(&article, id)
	c.JSON(http.StatusOK, gin.H{
		"message":  "删除成功",
		"redirect": "/articles/" + id, // 跳转URL
		"article":  article,           // 返回更新后的文章数据
	})
}
