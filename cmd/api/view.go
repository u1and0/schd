package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ViewList(c *gin.Context) {
	rows := data.Stack().Unstack()
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"a": rows,
	})
}
