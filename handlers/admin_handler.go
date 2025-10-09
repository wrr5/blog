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

func Admin(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "admin-index.html", gin.H{
		"user": user,
	})
}

func ShowAdminCategoriesPage(c *gin.Context) {
	var categories []models.Category

	size := config.AppConfig.Page.Size
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", size))
	offset := (page - 1) * pageSize

	var total int64
	// 先获取总记录数
	global.DB.Model(&models.Category{}).Count(&total)

	result := global.DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&categories)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "查询文章列表失败"})
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	Pagination := tools.Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
		BasePath:    "/admin/categories",
	}
	Pagination.CalculateDisplayPages(7)

	c.HTML(200, "admin-categories.html", gin.H{
		"CurrentPath": c.Request.URL.Path,
		"categories":  categories,
		"Pagination":  Pagination,
	})
}

func ShowAdminUsersPage(c *gin.Context) {
	c.HTML(200, "admin-users.html", gin.H{
		"CurrentPath": c.Request.URL.Path,
	})
}
