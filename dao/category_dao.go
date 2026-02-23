package dao

import (
	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
)

type CategoryDao struct{}

func (d *CategoryDao) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := global.DB.Find(&categories).Error
	return categories, err
}
