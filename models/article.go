package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Id       uint   `json:"id" gorm:"primarykey;autoIncrement:true"`
	Title    string `json:"title" gorm:"type:varchar(255);not null"`
	Content  string `json:"content" gorm:"type:text;not null"`
	IsPublic bool   `json:"is_public" gorm:"default:true"` // 是否公开

	// 添加外键关联 User
	UserID uint `json:"user_id" gorm:"not null;index"`                                                             // 外键字段
	User   User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 关联关系

	// 添加外键关联分类 (一篇文章属于一个分类)
	CategoryID *uint    `json:"category_id"` // 指针类型，允许 NULL
	Category   Category `json:"category" gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`

	// 关联标签 (一篇文章可以有多个标签)
	Tags []Tag `json:"tags,omitempty" gorm:"many2many:article_tags;" `

	// 关联评论，不推荐！
	// Comments []Comment `json:"comments" gorm:"foreignKey:ArticleID"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"` // 自动创建时间
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"` // 自动更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
