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
		// todo: 实现jwt验证进行登陆
		var user models.User
		result := global.DB.Where("email = ?", c.PostForm("email")).First(&user)
		if result.Error != nil || !tools.CheckPasswordHash(c.PostForm("password"), user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
			return
		}
		// 生成 JWT Token
		token, err := tools.GenerateJWT(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
			"user":    user,
			"token":   "Bearer " + token,
		})
	}

}
