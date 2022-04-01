package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
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
	c.HTML(http.StatusOK, "create.tmpl", "")
}

// CreateForm : Postするとフォームを読み取り実行
func CreateForm(c *gin.Context) {
	datum := ctrl.Datum{}
	if err := c.Bind(&datum); err != nil { // name, assign 取得
		c.Status(http.StatusBadRequest)
	}
	form := ctrl.Form{}
	if err := c.Bind(&form); err != nil { // id, noki-date, noki-misc 取得
		c.Status(http.StatusBadRequest)
	}
	id := ctrl.ID(form.ID0 + form.ID1)
	datum.Noki.Date = form.Date
	datum.Noki.Misc = form.Misc
	addData := ctrl.Data{id: datum}
	c.IndentedJSON(http.StatusOK, addData)
}

// Update : Put method
func Update(c *gin.Context) {
}

// Remove : Delete method
func Remove(c *gin.Context) {
}
