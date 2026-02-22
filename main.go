package main

import (
	"context"
	"fmt"

	"gitee.com/wwgzr/blog/config"
	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/router"
	"gitee.com/wwgzr/blog/tools"
	"github.com/thinkerou/favicon"
)

func main() {
	config.Init()
	tools.InitDB()
	rdb := global.RDB
	val, err := rdb.Get(context.Background(), "stuname").Result()
	if err == nil {
		fmt.Println("Value: " + val)
	} else {
		fmt.Println("Key not found")
	}

	r := router.SetupRouter()

	r.LoadHTMLGlob("templates/**/*.html")
	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")
	r.Use(favicon.New("./static/images/favicon.ico"))

	r.Run(":" + config.AppConfig.Server.Port)
}
