package middleware

import (
	"net/http"

	"gitee.com/wwgzr/blog/tools"
	"github.com/gin-gonic/gin"
)

// todo: JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Cookie 获取 token
		token, err := c.Cookie("auth_token")
		if err != nil {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		// 验证 token
		claims, err := tools.ParseJWT(token[7:])
		if err != nil {
			c.SetCookie("auth_token", "", -1, "/", "", false, true) // 清除无效 token
			// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		// 获取用户名
		username, ok := claims["username"].(string)
		if !ok {
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}

		// 将用户信息设置到上下文
		c.Set("username", username)
		c.Next()
	}
}
