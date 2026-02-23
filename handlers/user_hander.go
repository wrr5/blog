package handlers

import (
	"net/http"

	"gitee.com/wwgzr/blog/global"
	"gitee.com/wwgzr/blog/models"
	"gitee.com/wwgzr/blog/tools"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	global.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"users": users,
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
	var user models.User
	global.DB.First(&user, id)
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func UpdateUser(c *gin.Context) {
	type UpdateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
	}
	var jsdate UpdateUserRequest
	var user models.User
	if err := c.BindJSON(&jsdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	global.DB.First(&user, id)
	user.Username = jsdate.Username
	user.Password = jsdate.Password
	user.IsAdmin = jsdate.IsAdmin
	global.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功",
		"user":    user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	global.DB.Delete(&user, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}
