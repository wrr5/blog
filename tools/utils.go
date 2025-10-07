package tools

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var jwtSecret = []byte("my-secret-key")

// 生成 JWT Token
func GenerateJWT(userID uint, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(72 * time.Hour) // 72小时过期

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expireTime.Unix(),
	})
	return token.SignedString(jwtSecret)
}

// 解析JWT Token
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证 token 是否有效
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("无效的token")
	}
}

type Pagination struct {
	CurrentPage int
	TotalPages  int
	PageSize    int
	BasePath    string
}

// 在结构体上定义方法
func (p Pagination) NextPage() int {
	if p.CurrentPage >= p.TotalPages {
		return p.TotalPages
	}
	return p.CurrentPage + 1
}

func (p Pagination) PrevPage() int {
	if p.CurrentPage <= 1 {
		return 1
	}
	return p.CurrentPage - 1
}
