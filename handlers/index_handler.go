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
	offset := (page - 1) * pageSize

	var articles []models.Article
	var categories []models.Category
	global.DB.Find(&categories)
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
		BasePath:    "/",
	}
	Pagination.CalculateDisplayPages(7)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"user":       user,
		"title":      "文章列表",
		"articles":   articles,
		"categories": categories,
		"Pagination": Pagination,
	})
}
