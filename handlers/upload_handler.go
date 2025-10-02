package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 图片上传处理
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file[]")
	if err != nil {
		file, err = c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "上传失败: " + err.Error(),
			})
			return
		}
	}

	// 保存文件
	filename := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存文件失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "上传成功",
		"code": 0,
		"data": gin.H{
			"errFiles": []string{},
			"succMap": gin.H{
				file.Filename: "/" + filename,
			},
		},
	})
}
