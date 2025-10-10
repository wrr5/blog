package handlers

import (
	"fmt"
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
	uintID := uint(id)
	offset := (page - 1) * pageSize

	var articles []models.Article
	var categories []models.Category
	global.DB.Find(&categories)
	var total int64

	// 构建查询
	query := global.DB.Model(&models.Article{})

	// 如果提供了分类ID，添加分类过滤条件
	if uintID != 0 {
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
	Pagination := tools.Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
		BasePath:    "/",
	}
	Pagination.CalculateDisplayPages(7)
	fmt.Println(uintID)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"user":            user,
		"title":           "文章列表",
		"articles":        articles,
		"categories":      categories,
		"Pagination":      Pagination,
		"basePath":        "/",
		"currentCategory": uintID, // 传递当前选中的分类给模板
	})
}
