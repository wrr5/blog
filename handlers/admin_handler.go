package handlers

import (
	"net/http"

	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"github.com/gin-gonic/gin"
)

func Admin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_index.html", gin.H{})
}

func ShowAdminCategoriesPage(c *gin.Context) {
	var categories []models.Category
	global.DB.Find(&categories)

	c.HTML(200, "admin_categories.html", gin.H{
		"CurrentPath": c.Request.URL.Path,
		"categories":  categories,
	})
}

func ShowAdminUsersPage(c *gin.Context) {
	c.HTML(200, "admin_users.html", gin.H{
		"CurrentPath": c.Request.URL.Path,
	})
}
