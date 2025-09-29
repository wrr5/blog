package handlers

import (
	"net/http"

	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"gitee.com/wwgzr/blog/tools"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "用户列表",
	})
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if user.Password != c.PostForm("password2") {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "两次密码不一致！",
		})
		return
	}

	// 加密密码
	hashPassword, err := tools.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	user.Password = hashPassword

	// 新增用户
	result := global.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "注册成功",
		"redirect": "/auth/login",
		"user":     user,
	})
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"msg": "用户" + id + "信息",
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新用户" + id + "信息",
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除用户" + id + "信息",
	})
}
