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

func CreateCategory(c *gin.Context) {
	type CreateCategoryRequest struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	var jsdate CreateCategoryRequest
	var category models.Category
	if err := c.BindJSON(&jsdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category = models.Category{
		Name:        jsdate.Name,
		Slug:        jsdate.Slug,
		Description: jsdate.Description,
	}
	result := global.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "创建成功",
		"category": category, // 返回创建的类别数据
	})
}

func GetCategory(c *gin.Context) {
	id := c.Param("id")
	var caregory models.Category
	global.DB.First(&caregory, id)
	c.JSON(http.StatusOK, gin.H{
		"category": caregory,
	})
}

func UpdateCategory(c *gin.Context) {
	type UpdateCategoryRequest struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	var jsdate UpdateCategoryRequest
	var category models.Category
	if err := c.BindJSON(&jsdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	global.DB.First(&category, jsdate.ID)
	category.Name = jsdate.Name
	category.Slug = jsdate.Slug
	category.Description = jsdate.Description
	global.DB.Save(&category)
	c.JSON(http.StatusOK, gin.H{
		"message":  "修改成功",
		"category": category,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	global.DB.Delete(&category, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}
