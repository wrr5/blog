package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name        string `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Slug        string `json:"slug" gorm:"uniqueIndex;size:60;not null"` // URL友好的名称
	Description string `json:"description" gorm:"type:text"`

	// 添加外键关联 Article
	Articles []Article `json:"articles,omitempty" gorm:"many2many:article_tags;"` // 反向引用
}
