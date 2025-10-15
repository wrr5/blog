package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowAboutPage(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "about-index.html", gin.H{
		"user":       user,
		"CurrentURL": c.Request.URL.Path,
	})
}
