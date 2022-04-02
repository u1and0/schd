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

// Show : show 1 datum ID info
func Show(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	if datum, ok := data[id]; ok { // Cast
		c.HTML(http.StatusOK, "get.tmpl", gin.H{"id": id, "a": datum})
		return
	}
	c.HTML(http.StatusBadRequest, "get.tmpl", gin.H{"msg": "IDが見つかりません"})
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
	if err := addData.Add(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	// c.IndentedJSON(http.StatusOK, addData)
	c.HTML(http.StatusOK, "get.tmpl", gin.H{"id": id, "a": data[id]})
}

// Update : Put method
func Update(c *gin.Context) {
	c.HTML(http.StatusOK, "update.tmpl", "")
}

// Remove : Delete method
func Remove(c *gin.Context) {
}
