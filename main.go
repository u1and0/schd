package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ksn"
)

const (
	FILE = "test/sample.json"
	PORT = ":8080"
)

func readJSON(fs string) []byte {
	// Open file
	f, err := os.Open(fs)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	// Read data
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func main() {
	r := gin.Default()

	b := readJSON(FILE)
	data := ksn.Data{}
	json.Unmarshal(b, &data)
	fmt.Printf("%v", data)

	// API
	r.GET("/list", func(c *gin.Context) {
		rows := data.ToCalendar().ToRows()
		c.JSON(http.StatusOK, rows)
	})
	r.GET("/cal", func(c *gin.Context) {
		cal := data.ToCalendar()
		c.JSON(http.StatusOK, cal)
	})
	r.GET("/:id", func(c *gin.Context) {
		s := c.Param("id")
		id := ksn.ID(s) // Cast
		c.JSON(http.StatusOK, data[id])
	})
	r.GET("/all", func(c *gin.Context) {
		c.JSON(http.StatusOK, data)
	})

	r.Run(PORT)
}
