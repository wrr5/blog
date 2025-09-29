package tools

import (
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

var jwtSecret = []byte("your-secret-key")

// 生成 JWT Token
func GenerateJWT(username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(72 * time.Hour) // 72小时过期

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expireTime.Unix(),
	})
	return token.SignedString(jwtSecret)
}
