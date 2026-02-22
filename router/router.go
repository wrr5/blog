// router.go
package router

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"gitee.com/wwgzr/blog/handlers"
	"gitee.com/wwgzr/blog/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"iterate":          pageIterate,
		"buildCategoryURL": buildCategoryURL,
		"sliceContent":     sliceContent,
	})
	// 注册文章相关路由
	setIndexRoutes(r)
	setupArticleRoutes(r)
	setupUserRoutes(r)
	setupAuthRoutes(r)
	setupUploadRoutes(r)
	setupAdminRoutes(r)
	setupCategoryRoutes(r)
	setupCommentRoutes(r)
	r.GET("/search", middleware.AuthMiddleware(), handlers.ShowSearchResult)

	// 根路径跳转
	r.GET("/", middleware.OptionalAuthMiddleware, handlers.ShowIndex)
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

	// 使用用户认证中间件
	articleGroup.Use(middleware.AuthMiddleware())
	{
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

func setupCategoryRoutes(r *gin.Engine) {
	categoryGroup := r.Group("/categories")
	categoryGroup.Use(middleware.AuthMiddleware())
	{
		categoryGroup.GET("", handlers.GetCategories)
		categoryGroup.POST("", middleware.AdminRequired, handlers.CreateCategory)
		categoryGroup.GET("/:id", handlers.GetCategory)
		categoryGroup.PUT("/:id", middleware.AdminRequired, handlers.UpdateCategory)
		categoryGroup.DELETE("/:id", middleware.AdminRequired, handlers.DeleteCategory)
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

func setupUploadRoutes(r *gin.Engine) {
	// 文件上传路由组
	UploadGroup := r.Group("/upload")
	{
		UploadGroup.POST("", handlers.UploadFile)
	}
}

func setIndexRoutes(r *gin.Engine) {
	// 首页导航栏路由
	IndexGroup := r.Group("/")
	{
		IndexGroup.GET("home", middleware.AuthMiddleware(), handlers.ShowHomePage)
		IndexGroup.GET("about", middleware.OptionalAuthMiddleware, handlers.ShowAboutPage)
	}
}

func setupAdminRoutes(r *gin.Engine) {
	// 管理员路由组
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(), middleware.AdminRequired) // 先验证登录，再验证管理员权限
	{
		adminGroup.GET("", handlers.Admin)
		adminGroup.GET("/categories", handlers.ShowAdminCategoriesPage)
		adminGroup.GET("/users", handlers.ShowAdminUsersPage)
	}
}

func setupCommentRoutes(r *gin.Engine) {
	// 评论路由组
	adminGroup := r.Group("/comments")
	adminGroup.Use(middleware.AuthMiddleware())
	{
		adminGroup.GET("/:article_id", handlers.GetComment)
		adminGroup.POST("", handlers.CreateComment)
		adminGroup.DELETE("/:id", handlers.DeleteComment)
	}
}

func pageIterate(start, end int) []int {
	if start > end {
		return []int{}
	}
	items := make([]int, end-start+1)
	for i := range items {
		items[i] = start + i
	}
	return items
}

func buildCategoryURL(currentURL string, categoryID uint) string {
	u, err := url.Parse(currentURL)
	if err != nil {
		return ""
	}

	query := u.Query()
	query.Del("category")
	query.Set("category", strconv.FormatUint(uint64(categoryID), 10))
	query.Del("page")
	query.Del("pageSize")

	u.RawQuery = query.Encode()
	return u.String()
}

func sliceContent(content string, length int) string {
	// 移除Markdown格式的图片链接
	re := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	content = re.ReplaceAllString(content, "")

	// 移除HTML格式的图片标签
	re = regexp.MustCompile(`<img[^>]*>`)
	content = re.ReplaceAllString(content, "")

	// 移除普通的文件链接（以/uploads/开头的链接）
	re = regexp.MustCompile(`\[.*?\]\(/uploads/[^)]+\)`)
	content = re.ReplaceAllString(content, "")

	// 移除纯URL链接
	re = regexp.MustCompile(`https?://[^\s]+`)
	content = re.ReplaceAllString(content, "")

	// 清理多余的空格和换行
	content = strings.TrimSpace(content)
	re = regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")

	// 截取内容
	runes := []rune(content)
	if len(runes) > length {
		return string(runes[:length])
	}
	return content
}
