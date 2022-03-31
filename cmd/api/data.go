package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
)

const FILE = "test/sample.json"

var data = ctrl.Data{}

func init() {
	if err := data.ReadJSON(FILE); err != nil {
		panic(err)
	}
}

func Index(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data)
}

func Get(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	if datum, ok := data[id]; ok { // Cast
		c.IndentedJSON(http.StatusOK, datum)
		return
	}
	c.JSON(http.StatusBadRequest,
		gin.H{"error": fmt.Sprintf("%v not found", id)})
}

func Post(c *gin.Context) {
	var addData ctrl.Data
	if err := c.ShouldBindJSON(&addData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// Data exist check
	for k := range addData {
		if _, ok := data[k]; ok {
			msg := fmt.Sprintf("ID: %v データが既に存在しています。Updateを試してください。", k)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}
	// Set data
	for k, v := range addData {
		data[k] = v
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, addData)
}

func Delete(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	if _, ok := data[id]; !ok {
		msg := fmt.Sprintf("ID: %v が見つかりません。別のIDを指定してください。", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	delete(data, id)
	// Check deleted id
	if _, ok := data[id]; ok {
		msg := fmt.Sprintf("ID: %v を削除できませんでした。", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	msg := fmt.Sprintf("ID: %v を削除しました。", id)
	c.JSON(http.StatusOK, gin.H{"msg": msg})
}

func Put(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	var upData ctrl.Data
	if err := c.ShouldBindJSON(&upData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// Data exist check
	if _, ok := data[id]; !ok {
		msg := fmt.Sprintf("ID: %v データが存在しません。/postを試してください。", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	// Update data
	for k, v := range upData {
		data[k] = v
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, upData)
}

func List(c *gin.Context) {
	rows := data.Stack().Unstack()
	c.IndentedJSON(http.StatusOK, rows)
}