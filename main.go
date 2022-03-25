package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	FILE = "sample.json"
	PORT = ":8080"
)

type (
	Data  map[string]Datum
	Datum struct {
		Name   string `json:"機器名"`
		Assign string `json:"担当者"`
		Konpo  `json:"梱包"`
		Syuka  `json:"出荷日"`
		Noki   `json:"納期"`
	}
	Konpo struct {
		KonpoIrai  bool   `json:"梱包会社依頼要否"`
		WDH        string `json:"外寸法"`
		Mass       int    `json:"質量"`
		Yuso       string `json:"輸送手段"`
		Chaku      string `json:"到着予定日"`
		ToiawaseNo string `json:"問い合わせ番号"`
		Misc       string `json:"備考"`
	}
	Syuka struct {
		Date time.Time `json:"日付"`
		Misc string    `json:"備考"`
	}
	Noki struct {
		Date time.Time `json:"日付"`
		Misc string    `json:"備考"`
	}

	Rows []Row
	Row  struct {
		Konpo []Konpo
		Syuka []Syuka
		Noki  []Noki
	}
)

func (d *Data) toRows() ( rows Rows ) {
	for _, date := range dates(){
		for _,datum := range d{
			if datum.Konpo
		}
	}
	return 
}

func main() {
	// Open router
	r := gin.Default()

	// Open file
	jsonfile, err := os.Open(FILE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonfile.Close()

	// Data read
	b, err := ioutil.ReadAll(jsonfile)
	if err != nil {
		fmt.Println(err)
	}
	data := Data{}
	json.Unmarshal(b, &data)

	r.GET("/list", func(c *gin.Context) {
		rows := data.toRows()
		c.JSON(http.StatusOK, rows)
	})
	r.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, data[id])
	})
	r.GET("/all", func(c *gin.Context) {
		c.JSON(http.StatusOK, data)
	})

	r.Run(PORT)
}
