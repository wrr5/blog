package handlers

import (
	"net/http"

	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"gitee.com/wwgzr/blog/tools"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func Login(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	}
	if c.Request.Method == "POST" {
		var user models.User
		result := global.DB.Where("email = ?", c.PostForm("email")).First(&user)
		if result.Error != nil || !tools.CheckPasswordHash(c.PostForm("password"), user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
			return
		}
		// 生成 JWT Token
		token, err := tools.GenerateJWT(user.ID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
			return
		}

		// 设置 HTTP-only Cookie
		c.SetCookie("auth_token", "Bearer "+token, 3600*72, "/", "", false, true)
		// c.JSON(http.StatusOK, gin.H{
		// 	"message": "登录成功",
		// 	"user":    user,
		// })
		c.Redirect(http.StatusFound, "/")
	}
}

func Logout(c *gin.Context) {
	// 清除cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	// 返回成功信息或重定向到登录页
	// c.JSON(http.StatusOK, gin.H{"message": "退出成功"})
	c.Redirect(http.StatusFound, "/auth/login")
}
