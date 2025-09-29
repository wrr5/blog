// router.go
package router

import (
	"gitee.com/wwgzr/blog/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 注册文章相关路由
	setupArticleRoutes(r)
	setupUserRoutes(r)

	// 这里可以添加其他路由组，比如用户路由、认证路由等
	// setupAuthRoutes(r)

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "notfound.html", gin.H{"error": "页面不存在"})
	})

	return r
}

// setupArticleRoutes 配置文章相关路由
func setupArticleRoutes(r *gin.Engine) {
	articleHander := handlers.NewArticleHanders()

	// 文章路由组
	articleGroup := r.Group("/articles")
	{
		// 博客列表
		articleGroup.GET("", articleHander.ShowArticleList)
		// 发布博客页面
		articleGroup.GET("/new", articleHander.ShowCreateArticlePage)
		// 新增博客
		articleGroup.POST("", articleHander.CreateArticle)
		// 查询博客详情
		articleGroup.GET("/:id", articleHander.ShowArticleDetail)
		// 编辑博客页面
		articleGroup.GET("/edit/:id", articleHander.ShowArticleEdit)
		// 更新博客
		articleGroup.PUT("/:id", articleHander.UpdateArticle)
		// 删除博客
		articleGroup.DELETE("/:id", articleHander.DeleteArticle)
	}
}

func setupUserRoutes(r *gin.Engine) {
	// 用户路由组
	userGroup := r.Group("/users")
	{
		userGroup.POST("", handlers.Register)
	}

}
