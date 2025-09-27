package main

import (
	"gitee.com/wwgzr/blog/handlers"
	"gitee.com/wwgzr/blog/tools"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"gorm.io/gorm"
)

// 全局数据库变量
var DB *gorm.DB

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")
	r.Use(favicon.New("./static/1.ico"))

	DB = tools.InitDB("root", "123456")

	articleHander := handlers.NewArticleHanders(DB)
	homeHander := handlers.NewHomeHanders(DB)

	// 首页
	r.GET("/", homeHander.ShowHome)
	// 发布博客
	r.GET("/articles/new", articleHander.ShowCreateArticlePage)
	// 提交博客
	r.POST("/articles", articleHander.CreateArticle)

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "notfound.html", gin.H{"error": "页面不存在"})
	})

	r.Run(":8080")
}
