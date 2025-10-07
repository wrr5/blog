package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"gitee.com/wwgzr/blog/config"
	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"gitee.com/wwgzr/blog/tools"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// 文章处理组
type ArticleHanders struct {
}

func NewArticleHanders() *ArticleHanders {
	return &ArticleHanders{}
}

func (h *ArticleHanders) ShowArticleList(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "请先登录"})
		c.Abort()
		return
	}
	size := config.AppConfig.Page.Size
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", size))
	offset := (page - 1) * pageSize

	var articles []models.Article

	var total int64
	// 先获取总记录数
	global.DB.Model(&models.Article{}).Count(&total)

	result := global.DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "查询文章列表失败"})
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	Pagination := tools.Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
		BasePath:    "/articles",
	}
	Pagination.CalculateDisplayPages(7)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"user":       user,
		"title":      "文章列表",
		"articles":   articles,
		"Pagination": Pagination,
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
	type CreateArticleRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var newArticle CreateArticleRequest
	if err := c.BindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
			Title:   newArticle.Title,
			Content: newArticle.Content,
			User:    user,
		}

		result = global.DB.Create(&article)
		if result.Error != nil {
			c.HTML(http.StatusOK, "create.html", gin.H{
				"error": result.Error,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "发布成功",
			"article": article, // 返回创建的文章数据
		})
		// c.Redirect(http.StatusFound, "/articles")
	}
}

// Markdown 转 HTML
func mdToHTML(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func (h *ArticleHanders) ShowArticleDetail(c *gin.Context) {
	type temp struct {
		Id        uint
		Title     string
		Content   template.HTML
		UserID    uint
		UpdatedAt time.Time
		User      struct {
			Username string
		}
	}
	var article models.Article
	articleID := c.Param("id")
	global.DB.Preload("User").First(&article, articleID)
	article.Content = mdToHTML(article.Content)

	// 将markdown转换为HTML，并转换为template.HTML类型
	htmlContent := mdToHTML(article.Content)

	newArticle := temp{
		Id:        article.Id,
		Title:     article.Title,
		Content:   template.HTML(htmlContent),
		UserID:    article.UserID,
		UpdatedAt: article.UpdatedAt,
		User: struct {
			Username string
		}{
			Username: article.User.Username,
		},
	}
	c.HTML(http.StatusOK, "article.html", gin.H{
		"user_id": c.GetUint("user_id"),
		"article": newArticle,
	})
}

func (h *ArticleHanders) ShowArticleEdit(c *gin.Context) {
	// 获取文章id
	id := c.Param("id")
	var article models.Article
	global.DB.Preload("User").First(&article, id)

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
	})
}
