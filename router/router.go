// router.go
package router

import (
	"net/http"

	"gitee.com/wwgzr/blog/handlers"
	"gitee.com/wwgzr/blog/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 注册文章相关路由
	setupArticleRoutes(r)
	setupUserRoutes(r)
	setupAuthRoutes(r)

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "notfound.html", gin.H{"error": "页面不存在"})
	})

	return r
}

// setupArticleRoutes 配置文章相关路由
func setupArticleRoutes(r *gin.Engine) {
	articleHander := handlers.NewArticleHanders()

	// 根路径跳转到, 已登录到文章列表, 未登录到登陆页
	r.GET("/", func(c *gin.Context) {
		// todo: 未登录到登陆页面
		c.Redirect(http.StatusFound, "/articles")
	})
	// 文章路由组
	articleGroup := r.Group("/articles")

	// 使用用户认证中间件
	articleGroup.Use(middleware.AuthMiddleware())
	{
		// 博客列表
		articleGroup.GET("", articleHander.ShowArticleList)
		// 发布博客页面
		articleGroup.GET("/create", articleHander.ShowCreateArticlePage)
		// 新增博客
		articleGroup.POST("", articleHander.CreateArticle)
		// 查询博客详情
		articleGroup.GET("/:id", articleHander.ShowArticleDetail)
		// 编辑博客页面
		articleGroup.GET("/:id/edit", articleHander.ShowArticleEdit)
		// 更新博客
		articleGroup.PUT("/:id", articleHander.UpdateArticle)
		// 删除博客
		articleGroup.DELETE("/:id", articleHander.DeleteArticle)
	}
}

func setupUserRoutes(r *gin.Engine) {
	// 用户资源路由组
	userGroup := r.Group("/users")

	// 使用用户认证中间件
	userGroup.Use(middleware.AuthMiddleware())
	{
		// 获取用户列表 - GET /users
		userGroup.GET("", handlers.GetUsers)

		// 创建用户（注册）- POST /users
		userGroup.POST("", handlers.CreateUser)

		// 获取特定用户 - GET /users/:id
		userGroup.GET("/:id", handlers.GetUser)

		// 更新用户 - PUT /users/:id
		userGroup.PUT("/:id", handlers.UpdateUser)

		// 删除用户 - DELETE /users/:id
		userGroup.DELETE("/:id", handlers.DeleteUser)
	}

}

func setupAuthRoutes(r *gin.Engine) {
	// 认证路由组
	authGroup := r.Group("/auth")
	{
		// 显示登录页面 - GET /auth/login
		authGroup.GET("/login", handlers.Login)

		// 执行登录 - POST /auth/login
		authGroup.POST("/login", handlers.Login)

		// 退出登录 - GET /auth/logout
		authGroup.GET("/logout", middleware.AuthMiddleware(), handlers.Logout)

		// 显示注册页面 - GET /auth/register
		authGroup.GET("/register", handlers.Register)

		// 注册账号 - POST /auth/register
		authGroup.POST("/register", handlers.CreateUser)

		// 登出 - POST /auth/logout
		// authGroup.POST("/logout", handlers.Logout)
	}
}
