package api

import (
	"fmt"
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
	if datum, ok := data[id]; ok {
		c.HTML(http.StatusOK, "get.tmpl", gin.H{"id": id, "a": datum})
		return
	}
	c.HTML(http.StatusBadRequest, "get.tmpl", gin.H{"msg": "IDが見つかりません"})
}

// CreateForm : Post
func CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create.tmpl", "")
}

// Create : Postするとフォームを読み取り実行
func Create(c *gin.Context) {
	datum := ctrl.Datum{}
	if err := c.Bind(&datum); err != nil { // name, assign 取得
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
		return
	}
	form := ctrl.Form{}
	if err := c.Bind(&form); err != nil { // id, noki-date, noki-misc 取得
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
		return
	}
	id := ctrl.ID(form.ID0 + form.ID1)
	datum.Noki.Date = form.Noki.Date
	datum.Noki.Misc = form.Noki.Misc
	addData := ctrl.Data{id: datum}
	if err := addData.Add(&data); err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
		return
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	c.HTML(http.StatusOK, "get.tmpl", gin.H{"id": id, "a": data[id]})
}

// RefreshForm : 更新情報を入力するフォーム GETされてHTMLを返す
func RefreshForm(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	if datum, ok := data[id]; ok {
		c.HTML(http.StatusOK, "update.tmpl", gin.H{"id": id, "a": datum})
		return
	}
	c.HTML(http.StatusBadRequest, "update.tmpl", gin.H{"msg": "IDが見つかりません"})
}

// Refresh : Post from RefreshForm
func Refresh(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	upData := ctrl.Datum{}
	if err := c.Bind(&upData); err != nil { // name, assign 取得
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
		return
	}
	// form := ctrl.Form{}
	// if err := c.Bind(&form); err != nil { // id, noki-date, noki-misc 取得
	// 	c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
	// 	return
	// }
	// konpo := ctrl.Konpo{}
	fmt.Printf("Datamu %#v\n", upData)
	if err := c.Bind(&upData.Konpo); err != nil { // id, noki-date, noki-misc 取得
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
		return
	}
	fmt.Printf("Konpo %#v\n", upData)
	if err := c.Bind(&upData.Syuka); err != nil { // id, noki-date, noki-misc 取得
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
		return
	}
	fmt.Printf("Syua %#v\n", upData)
	if err := c.Bind(&upData.Noki); err != nil { // id, noki-date, noki-misc 取得
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
		return
	}
	fmt.Printf("Noki %#v\n", upData)
	if err := upData.Update(id, &data); err != nil {
		// c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
		c.JSON(http.StatusBadRequest, upData)
		return
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	c.HTML(http.StatusOK, "get.tmpl", gin.H{"a": data[id]})
}

// Remove : Delete method
func Remove(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	if err := id.Del(&data); err != nil {
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"error": err})
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	c.HTML(http.StatusOK, "delete.tmpl", gin.H{"id": c.Param("id")})
}
