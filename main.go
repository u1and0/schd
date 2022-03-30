package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ksn"
)

const (
	FILE = "test/sample.json"
	PORT = ":8080"
)

func main() {
	// Read data
	b := ksn.ReadJSON(FILE)
	data := ksn.Data{}
	json.Unmarshal(b, &data)

	// Router
	r := gin.Default()
	r.LoadHTMLGlob("template/*.tmpl")

	// API
	v1 := r.Group("api/v1")
	{
		d := v1.Group("/data")
		{
			d.GET("/all", func(c *gin.Context) {
				c.IndentedJSON(http.StatusOK, data)
			})
			d.GET("/:id", func(c *gin.Context) {
				s := c.Param("id")
				id := ksn.ID(s)
				if datum, ok := data[id]; ok { // Cast
					c.IndentedJSON(http.StatusOK, datum)
					return
				}
				c.JSON(http.StatusBadRequest,
					gin.H{"msg": fmt.Sprintf("%v not found", id)})
			})
			d.GET("/cal", func(c *gin.Context) {
				cal := data.Stack()
				c.IndentedJSON(http.StatusOK, cal)
			})
			d.GET("/list", func(c *gin.Context) {
				rows := data.Stack().Unstack()
				c.IndentedJSON(http.StatusOK, rows)
			})
			d.POST("/post", func(c *gin.Context) {
				var addData ksn.Data
				if err := c.ShouldBindJSON(&addData); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
				}
				// Data exist check
				for k := range addData {
					if _, ok := data[k]; ok {
						msg := fmt.Sprintf("ID: %v データが既に存在しています。Updateを試してください。", k)
						c.JSON(http.StatusBadRequest, gin.H{"msg": msg})
						return
					}
				}
				// Append data
				for k, v := range addData {
					data[k] = v
				}
				c.IndentedJSON(http.StatusOK, addData)
			})
		}

		v := v1.Group("view")
		{
			v.GET("/list", func(c *gin.Context) {
				rows := data.Stack().Unstack()
				c.HTML(http.StatusOK, "index.tmpl", gin.H{
					"a": rows,
				})
			})
		}
	}

	r.Run(PORT)
}
