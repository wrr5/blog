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

type CreateUserRequest struct {
	models.User        // 嵌入 User 结构体，包含 Username, Email, Password, IsAdmin 等字段
	Password2   string `form:"password2" json:"password2" binding:"required"` // 确认密码
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	// ShouldBind 根据 Content-Type 自动解析 JSON 或表单
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证密码一致性
	if req.Password != req.Password2 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "两次密码不一致！"})
		return
	}

	// 加密密码
	hashPassword, err := tools.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.User.Password = hashPassword

	// 创建用户
	result := global.DB.Create(&req.User)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "注册成功",
		"redirect": "/auth/login",
		"user":     req.User,
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
		Email    string `json:"email"`
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
	user.Email = jsdate.Email
	// 加密密码
	hashPassword, err := tools.HashPassword(jsdate.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = hashPassword
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
