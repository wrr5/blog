package handlers

import (
	"net/http"
	"strconv"

	"gitee.com/wwgzr/blog/config"
	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"gitee.com/wwgzr/blog/tools"
	"github.com/gin-gonic/gin"
)

func ShowIndex(c *gin.Context) {
	user, _ := c.Get("user")
	// if !exists {
	// 	c.JSON(401, gin.H{"error": "请先登录"})
	// 	c.Abort()
	// 	return
	// }
	size := config.AppConfig.Page.Size
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", size))
	categoryID := c.Query("category") // 获取分类参数
	id, _ := strconv.ParseUint(categoryID, 10, 0)
	categoryUintID := uint(id)
	offset := (page - 1) * pageSize

	var articles []models.Article
	var categories []models.Category
	global.DB.Find(&categories)
	var total int64

	// 构建查询
	query := global.DB.Model(&models.Article{}).Where("is_public = ?", true)

	// 如果提供了分类ID，添加分类过滤条件
	if categoryUintID != 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 先获取总记录数（带分类条件）
	query.Count(&total)

	// 查询文章列表（带分类条件）
	result := query.Preload("User").Preload("Category").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "查询文章列表失败"})
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	// 获取当前请求的所有查询参数
	queryParams := c.Request.URL.Query()

	// 移除分页相关的参数，避免重复
	queryParams.Del("page")
	queryParams.Del("pageSize")

	// 构建基础路径，包含现有查询参数
	basePath := c.Request.URL.Path + "?" + queryParams.Encode()
	if len(queryParams) > 0 {
		basePath += "&"
	}
	Pagination := tools.Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
		BasePath:    basePath,
	}
	Pagination.CalculateDisplayPages(7)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"user":            user,
		"title":           "首页",
		"articles":        articles,
		"categories":      categories,
		"Pagination":      Pagination,
		"currentCategory": categoryUintID, // 传递当前选中的分类给模板
		"BaseURL":         "/",
		"CurrentURL":      c.Request.URL.Path,
	})
}
