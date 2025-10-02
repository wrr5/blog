package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 图片上传处理 - 支持多个文件
func UploadFile(c *gin.Context) {
	// 获取多文件表单
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "上传失败: " + err.Error(),
		})
		return
	}

	// 获取所有上传的文件
	files := form.File["file[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "没有找到上传的文件",
		})
		return
	}

	// 创建上传目录
	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建目录失败: " + err.Error(),
		})
		return
	}

	succMap := make(gin.H)
	errFiles := []string{}

	// 遍历处理每个文件
	for _, file := range files {
		// 生成安全文件名（避免覆盖和路径遍历攻击）
		filename := "uploads/" + generateSafeFilename(file.Filename)

		// 保存文件
		if err := c.SaveUploadedFile(file, filename); err != nil {
			errFiles = append(errFiles, file.Filename)
			continue
		}

		succMap[file.Filename] = "/" + filename
	}

	// 如果所有文件都上传失败
	if len(succMap) == 0 && len(errFiles) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "所有文件保存失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  fmt.Sprintf("上传完成，成功%d个，失败%d个", len(succMap), len(errFiles)),
		"code": 0,
		"data": gin.H{
			"errFiles": errFiles,
			"succMap":  succMap,
		},
	})
}

// 生成安全文件名
func generateSafeFilename(original string) string {
	// 提取文件扩展名
	ext := filepath.Ext(original)
	// 生成随机文件名
	name := fmt.Sprintf("%d_%s", time.Now().UnixNano(), uuid.New().String())
	return name + ext
}
