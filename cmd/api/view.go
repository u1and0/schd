package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// View : 日付リスト
func View(c *gin.Context) {
	rows := data.Stack().Unstack()
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"a": rows,
	})
}

// Create : Post
func Create(c *gin.Context) {
	c.HTML(http.StatusOK, "create.tmpl", gin.H{"a": data})
}

// Update : Put method
func Update(c *gin.Context) {
}

// Remove : Delete method
func Remove(c *gin.Context) {
}
