package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Id      uint   `json:"id" gorm:"primarykey;autoIncrement:true"`
	Title   string `json:"title" gorm:"type:varchar(255);not null"`
	Content string `json:"content" gorm:"type:text;not null"`

	// 添加外键关联 User
	UserID uint `json:"user_id" gorm:"not null;index"`                                                             // 外键字段
	User   User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 关联关系

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"` // 自动创建时间
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"` // 自动更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
