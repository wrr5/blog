package main

import (
	"gitee.com/wwgzr/blog/config"
	"gitee.com/wwgzr/blog/router"
	"gitee.com/wwgzr/blog/tools"
	"github.com/thinkerou/favicon"
)

func main() {
	config.Init()
	tools.InitDB()

	r := router.SetupRouter()

	r.LoadHTMLGlob("templates/**/*.html")
	r.Static("/static", "./static")
	r.Use(favicon.New("./static/1.ico"))

	r.Run(":" + config.AppConfig.Server.Port)
}
