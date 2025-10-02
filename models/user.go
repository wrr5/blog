package models

import (
	"gorm.io/gorm"
)

type User struct {
	// gorm.Model 的定义
	// ID        uint           `gorm:"primaryKey;autoIncrement:true"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	gorm.Model
	Username string    `json:"username" form:"username" gorm:"uniqueIndex;size:50;not null"`
	Email    string    `json:"email" form:"email" gorm:"uniqueIndex;size:100;not null"`
	Password string    `json:"password" form:"password" gorm:"size:255;not null"`
	Articles []Article `json:"articles" gorm:"foreignKey:UserID"` // 反向引用
}
