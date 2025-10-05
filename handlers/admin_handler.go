package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Admin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_index.html", gin.H{})
}

func ShowAdminCategoriesPage(c *gin.Context) {
	c.HTML(200, "admin_categories.html", gin.H{
		"CurrentPath": c.Request.URL.Path,
	})
}

func ShowAdminUsersPage(c *gin.Context) {
	c.HTML(200, "admin_users.html", gin.H{
		"CurrentPath": c.Request.URL.Path,
	})
}
