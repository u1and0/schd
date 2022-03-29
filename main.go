package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ksn"
)

const (
	FILE = "test/sample.json"
	PORT = ":8080"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("static/*.tmpl")

	// Read data
	b := ksn.ReadJSON(FILE)
	data := ksn.Data{}
	json.Unmarshal(b, &data)

	// API
	r.GET("/all", func(c *gin.Context) {
		c.JSON(http.StatusOK, data)
	})
	r.GET("/:id", func(c *gin.Context) {
		s := c.Param("id")
		id := ksn.ID(s) // Cast
		c.JSON(http.StatusOK, data[id])
	})
	r.GET("/cal", func(c *gin.Context) {
		cal := data.ToCalendar()
		c.JSON(http.StatusOK, cal)
	})
	r.GET("/list-json", func(c *gin.Context) {
		rows := data.ToCalendar().ToRows()
		c.JSON(http.StatusOK, rows)
	})
	r.GET("/list", func(c *gin.Context) {
		rows := data.ToCalendar().ToRows()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"a": rows,
		})
	})

	r.Run(PORT)
}
