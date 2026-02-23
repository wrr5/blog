package dao

import (
	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
)

type ArticleDao struct{}

func (d *ArticleDao) GetList(userID, categoryID uint, offset, limit int) ([]models.Article, error) {
	var articles []models.Article
	query := global.DB.Model(&models.Article{}).
		Where("(is_public = ? OR (is_public = ? AND user_id = ?))", true, false, userID)
	if categoryID != 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	err := query.Preload("User").Preload("Category").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&articles).Error
	return articles, err
}

func (d *ArticleDao) Count(userID, categoryID uint) (int64, error) {
	var total int64
	query := global.DB.Model(&models.Article{}).
		Where("(is_public = ? OR (is_public = ? AND user_id = ?))", true, false, userID)
	if categoryID != 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	err := query.Count(&total).Error
	return total, err
}
