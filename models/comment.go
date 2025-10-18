package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ArticleID uint   `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	ParentID  *uint  // 指向父评论的ID，如果为nil则表示是顶级评论
	Content   string `gorm:"type:text;not null"`

	User    User      `gorm:"foreignKey:UserID"`
	Article Article   `gorm:"foreignKey:ArticleID"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"-"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// 创建评论请求结构
type CreateCommentRequest struct {
	Content   string `json:"content" binding:"required"`
	ArticleID string `json:"article_id" binding:"required"`
	ParentID  string `json:"parent_id"` // 可以是空字符串或者数字字符串
}

// 评论响应结构
type CommentResponse struct {
	ID        uint              `json:"id"`
	Content   string            `json:"content"`
	CreatedAt time.Time         `json:"created_at"`
	User      UserResponse      `json:"user"`
	Replies   []CommentResponse `json:"replies,omitempty"`
}

// 用户响应结构（简化版，避免敏感信息）
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
