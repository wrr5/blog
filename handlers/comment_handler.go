// handlers/comment.go
package handlers

import (
	"net/http"
	"strconv"

	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetComment 获取文章评论列表
func GetComment(c *gin.Context) {
	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的文章ID",
		})
		return
	}

	// 检查文章是否存在
	var article models.Article
	if err := global.DB.First(&article, uint(articleID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "文章不存在",
		})
		return
	}

	// 获取顶级评论（parent_id为NULL的评论）
	var comments []models.Comment
	if err := global.DB.
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("created_at ASC")
		}).
		Where("article_id = ? AND parent_id IS NULL", uint(articleID)).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取评论失败",
		})
		return
	}

	// 转换为响应格式
	commentResponses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = convertToCommentResponse(comment)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "获取评论成功",
		"comments": commentResponses,
	})
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从上下文中获取用户ID（AuthMiddleware设置的）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未认证",
		})
		return
	}

	// 转换字符串为 uint
	articleID, err := strconv.ParseUint(req.ArticleID, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的文章ID格式",
		})
		return
	}

	// 检查文章是否存在
	var article models.Article
	if err := global.DB.First(&article, uint(articleID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "文章不存在",
		})
		return
	}

	parentID, _ := strconv.ParseUint(req.ParentID, 10, 32)
	// 如果parent_id不为空，检查父评论是否存在且属于同一文章
	if parentID != 0 {
		var parentComment models.Comment
		if err := global.DB.First(&parentComment, parentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "父评论不存在",
			})
			return
		}
		if parentComment.ArticleID != uint(articleID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "父评论不属于该文章",
			})
			return
		}
	}

	// 创建评论
	parentIDUint := uint(parentID)
	comment := models.Comment{
		ArticleID: uint(articleID),
		UserID:    userID.(uint),
		ParentID:  &parentIDUint,
		Content:   req.Content,
	}
	if parentIDUint == 0 {
		comment.ParentID = nil
	}

	if err := global.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建评论失败",
		})
		return
	}

	// 重新加载评论以获取关联的用户信息
	if err := global.DB.Preload("User").First(&comment, comment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取评论详情失败",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "评论创建成功",
		"comment": convertToCommentResponse(comment),
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	// 获取评论ID
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的评论ID",
		})
		return
	}

	// 获取当前用户
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "请先登录",
		})
		return
	}

	// 查询评论信息
	var comment models.Comment
	if err := global.DB.Preload("Article").First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "评论不存在",
		})
		return
	}

	// 检查权限：文章作者或评论发布者
	if comment.UserID != userID.(uint) && comment.Article.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权删除此评论",
		})
		return
	}

	// 开启事务
	tx := global.DB.Begin()

	// 如果是一级评论，删除所有子评论
	if comment.ParentID == nil {
		if err := tx.Where("parent_id = ?", comment.ID).Delete(&models.Comment{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "删除子评论失败",
			})
			return
		}
	}

	// 删除评论本身
	if err := tx.Delete(&comment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除评论失败",
		})
		return
	}

	// 提交事务
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// convertToCommentResponse 将Comment转换为CommentResponse
func convertToCommentResponse(comment models.Comment) models.CommentResponse {
	response := models.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		User: models.UserResponse{
			ID:       comment.User.ID,
			Username: comment.User.Username,
		},
	}

	// 递归转换回复
	if len(comment.Replies) > 0 {
		response.Replies = make([]models.CommentResponse, len(comment.Replies))
		for i, reply := range comment.Replies {
			response.Replies[i] = convertToCommentResponse(reply)
		}
	}

	return response
}
