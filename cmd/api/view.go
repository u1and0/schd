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

// Create : Post
func Create(c *gin.Context) {
	c.HTML(http.StatusOK, "create.tmpl", "")
}

func CreateForm(c *gin.Context) {
	id0 := c.PostForm("id0")
	id1 := c.PostForm("id1")
	id := ctrl.ID(id0 + id1)
	name := c.PostForm("name")
	assign := c.PostForm("assign")
	ds := c.PostForm("date")
	// y, _ := strconv.Atoi(ds[:4])
	// m, _ := strconv.Atoi(ds[:4])
	// d, _ := strconv.Atoi(ds[:4])
	// date := time.Date(strconv.Atoi)
	misc := c.PostForm("misc")
	// addData := ctrl.Data{}
	noki := ctrl.Noki{Misc: misc}
	datum := ctrl.Datum{Name: name, Assign: assign, Noki: noki}
	// addData[id] = datum
	fmt.Printf("datum: %#v\n", datum)
	fmt.Printf("date: %v\n", ds)
	fmt.Printf("date: %v\n", id)
}

// Update : Put method
func Update(c *gin.Context) {
}

// Remove : Delete method
func Remove(c *gin.Context) {
}
