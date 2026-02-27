package service

import (
	"gitee.com/wwgzr/blog/dao"
	"gitee.com/wwgzr/blog/models"
)

// type ArticleService interface {
//     GetHomePageData(userID, categoryID uint, page, pageSize int) ([]models.Article, int64, []models.Category, error)
// }

type ArticleService struct {
	articleDao  *dao.ArticleDao
	categoryDao *dao.CategoryDao
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		articleDao:  &dao.ArticleDao{},
		categoryDao: &dao.CategoryDao{},
	}
}

// GetHomePageData 获取首页所需数据：文章列表、总数、分类列表
func (s *ArticleService) GetHomePageData(userID, categoryID uint, page, pageSize int) (articles []models.Article, total int64, categories []models.Category, err error) {
	// 获取分类列表（也可单独调用，但这里一起返回方便）
	categories, err = s.categoryDao.GetAll()
	if err != nil {
		return
	}

	offset := (page - 1) * pageSize
	articles, err = s.articleDao.GetList(userID, categoryID, offset, pageSize)
	if err != nil {
		return
	}

	total, err = s.articleDao.Count(userID, categoryID)
	return
}
