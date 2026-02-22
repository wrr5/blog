package handlers

import (
	"fmt"
	"net/http"

	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"github.com/gin-gonic/gin"
)

func ShowSearchResult(c *gin.Context) {
	// 1. 获取当前登录用户（从中间件存入的上下文）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	currentUser := user.(models.User) // 假设你的 User 模型在 models 包中

	// 2. 获取搜索关键词
	keyword := c.Query("q")
	if keyword == "" {
		// 关键词为空，可以返回空列表或提示
		c.JSON(http.StatusOK, gin.H{
			"user":     currentUser,
			"articles": []models.Article{},
		})
		return
	}

	// 3. 使用 GORM 进行模糊查询
	var articles []models.Article
	db := global.DB

	// 构建基础查询：只搜索当前用户有权看到的文章
	// 条件：is_public = true OR (user_id = 当前用户ID)
	// 并且标题或内容包含关键词
	err := db.Where("is_public = ? OR user_id = ?", true, currentUser.ID).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		// 预加载关联信息，方便前端展示
		Preload("User").
		Preload("Category").
		Preload("Tags").
		// 按更新时间倒序排列（可选）
		Order("updated_at DESC").
		Find(&articles).Error

	if err != nil {
		// 查询出错
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}

	// 4. 返回 JSON 结果
	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("共查询到%d条记录", len(articles)),
		"count":    len(articles),
		"articles": articles,
	})
}
