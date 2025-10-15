package handlers

import (
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

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

func (h *ArticleHanders) ShowCreateArticlePage(c *gin.Context) {
	var categories []models.Category
	global.DB.Find(&categories)

	if username, ok := c.Get("username"); ok {
		c.HTML(http.StatusOK, "create.html", gin.H{
			"username":   username,
			"categories": categories,
		})
	}
}

func (h *ArticleHanders) CreateArticle(c *gin.Context) {
	type CreateArticleRequest struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryId string `json:"category_id"`
		IsPublic   string `json:"is_public"`
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
	}

	var article models.Article
	// 创建文章对象
	var categoryID *uint
	if newArticle.CategoryId != "0" {
		id, err := strconv.ParseUint(newArticle.CategoryId, 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "无效的分类ID"})
			return
		}
		idUint := uint(id)
		categoryID = &idUint
	}
	article = models.Article{
		Title:      newArticle.Title,
		Content:    newArticle.Content,
		User:       user,
		CategoryID: categoryID,
	}

	result := global.DB.Create(&article)
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
		IsPublic  bool
		User      struct {
			Username string
		}
		Category struct {
			ID   uint
			Name string
		}
	}
	var article models.Article
	articleID := c.Param("id")
	global.DB.Preload("Category").Preload("User").First(&article, articleID)

	if !tools.VisitPrivateArticle(c, article) {
		return
	}
	article.Content = mdToHTML(article.Content)

	// 将markdown转换为HTML，并转换为template.HTML类型
	htmlContent := mdToHTML(article.Content)

	newArticle := temp{
		Id:        article.Id,
		Title:     article.Title,
		Content:   template.HTML(htmlContent),
		UserID:    article.UserID,
		UpdatedAt: article.UpdatedAt,
		IsPublic:  article.IsPublic,
		User: struct {
			Username string
		}{
			Username: article.User.Username,
		},
		Category: struct {
			ID   uint
			Name string
		}{
			ID:   article.Category.ID,
			Name: article.Category.Name,
		},
	}
	// 判断是从首页还是从我的文章跳转支持
	referer := c.Request.Header.Get("Referer")
	lastPart := path.Base(referer)
	var BaseURL string
	// 根路径是"home"是从我的文章跳转至此
	if strings.HasPrefix(lastPart, "home") {
		BaseURL = "home"
	}
	c.HTML(http.StatusOK, "article.html", gin.H{
		"user_id": c.GetUint("user_id"),
		"article": newArticle,
		"BaseURL": BaseURL,
	})
}

func (h *ArticleHanders) ShowArticleEdit(c *gin.Context) {
	// 获取文章id
	id := c.Param("id")
	var article models.Article
	global.DB.Preload("User").Preload("Category").First(&article, id)

	var categories []models.Category
	global.DB.Find(&categories)

	if !tools.VisitPrivateArticle(c, article) {
		return
	}
	// 只有自己可以编辑
	userID := c.GetUint("user_id")
	if article.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "无权操作此文章"})
		c.Abort()
		return
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"article":    article,
		"categories": categories,
	})
}

func (h *ArticleHanders) UpdateArticle(c *gin.Context) {
	type UpdateArticleRequest struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryId string `json:"category_id"`
		IsPublic   string `json:"is_public"`
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

	// 创建文章对象
	var categoryID *uint
	if req.CategoryId != "0" {
		id, err := strconv.ParseUint(req.CategoryId, 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "无效的分类ID"})
			return
		}
		idUint := uint(id)
		categoryID = &idUint
	}
	// 修改文章
	article.Title = req.Title
	article.Content = req.Content
	article.CategoryID = categoryID
	switch req.IsPublic {
	case "0":
		article.IsPublic = false
	case "1":
		article.IsPublic = true
	}
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
