package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Id        uint           `json:"id" gorm:"primarykey;autoIncrement:true"`
	Title     string         `json:"title" gorm:"type:varchar(255);not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"` // 自动创建时间
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"` // 自动更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
