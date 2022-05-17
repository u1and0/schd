package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
)

// FILE : DB file path
const (
	FILE        = "test/sample.json"
	ADDRESSFILE = "db/住所録.json"
	SECTIONFILE = "db/課.json"
)

var data = ctrl.Data{}

type (
	// AddressMap : JSON ファイルから読み取った住所録
	AddressMap map[string]Address
	// Address : 1住所あたり5行まで, 1行あたり15文字まで
	// Excelシートの枠の都合
	Address []string
)

func init() {
	if err := data.ReadJSON(FILE); err != nil {
		// panic にするとgo test時にFILEが見つからないエラー
		fmt.Println(err)
	}
}

// All : show all data
func All(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data)
}

// Get : show 1 datum by id
func Get(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	if datum, ok := data[id]; ok { // Cast
		c.IndentedJSON(http.StatusOK, datum)
		return
	}
	c.JSON(http.StatusBadRequest,
		gin.H{"error": fmt.Sprintf("%v not found", id)})
}

// Post : Create some data from JSON
func Post(c *gin.Context) {
	var addData ctrl.Data
	if err := c.ShouldBindJSON(&addData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := addData.Add(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, addData)
}

// Delete : Delete 1 datum by id
func Delete(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	if err := id.Del(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	c.JSON(http.StatusNoContent, gin.H{"id": id})
}

// Put : Update 1 datum by ID
func Put(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	var upData ctrl.Datum
	if err := c.ShouldBindJSON(&upData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := upData.Update(id, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := data.WriteJSON(FILE); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, id)
}

// List : Show table like by date
func List(c *gin.Context) {
	rows := data.Stack().Unstack().Verbose(data)
	c.IndentedJSON(http.StatusOK, rows)
}

// FetchAddress : 住所録をJSONで返す
func FetchAddress(c *gin.Context) {
	var m AddressMap
	if err := ctrl.UnmarshalJSONfile(m, ADDRESSFILE); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, m)
}

func (a *Address) String() string {
	return strings.Join(*a, "\n")
}
