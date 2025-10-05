package handlers

import (
	"net/http"

	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category
	global.DB.Find(&categories)
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
