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
		Date       time.Time `json:"日付"`
		KonpoIrai  bool      `json:"梱包会社依頼要否"`
		WDH        string    `json:"外寸法"`
		Mass       int       `json:"質量"`
		Yuso       string    `json:"輸送手段"`
		Chaku      string    `json:"到着予定日"`
		ToiawaseNo string    `json:"問合わせ番号"`
		Misc       string    `json:"備考"`
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
		Date        time.Time `json:"日付"`
		KonpoID     string    `json:"梱包-生産番号"`
		KonpoName   string    `json:"梱包-機器名"`
		KonpoAssign string    `json:"梱包-担当者"`
		KonpoIrai   bool      `json:"梱包会社依頼要否"`
		WDH         string    `json:"外寸法"`
		Mass        int       `json:"質量"`
		Yuso        string    `json:"輸送手段"`
		Chaku       string    `json:"到着予定日"`
		ToiawaseNo  string    `json:"問合わせ番号"`
		KonpoMisc   string    `json:"梱包-備考"`

		SyukaID     string `json:"出荷-生産番号"`
		SyukaName   string `json:"出荷-機器名"`
		SyukaAssign string `json:"出荷-担当者"`
		SyukaMisc   string `json:"出荷-備考"`

		NokiID     string `json:"納期-生産番号"`
		NokiName   string `json:"納期-機器名"`
		NokiAssign string `json:"納期-担当者"`
		NokiMisc   string `json:"納期-備考"`
	}
)

// ToCalendar : Rowsのテーブルを返す
// JavaScriptでHTML テーブル
func (d *Data) ToCalendar() (rows Rows) {
	dates := []time.Time{
		time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 24, 0, 0, 0, 0, time.UTC),
	}
	for _, date := range dates {
		row := Row{}
		for id, datum := range *d {
			switch date {
			case datum.Konpo.Date:
				fmt.Printf("Konpo date 追加: %s,%s\n", id, date)
				row.KonpoID = id
				row.KonpoName = datum.Name
				row.KonpoAssign = datum.Assign
				row.WDH = datum.WDH
				fallthrough
			case datum.Syuka.Date:
				fmt.Printf("Syuka date 追加: %s,%s\n", id, date)
				row.SyukaID = id
				row.SyukaName = datum.Name
				row.SyukaAssign = datum.Assign
				fallthrough
			case datum.Noki.Date:
				fmt.Printf("Noki date 追加: %s,%s\n", id, date)
				row.NokiID = id
				row.NokiName = datum.Name
				row.NokiAssign = datum.Assign
				fallthrough
			default:
				row.Date = date
			}
			rows = append(rows, row)
		}
	}
	return
}

func main() {
	r := gin.Default()

	// Open file
	jsonfile, err := os.Open(FILE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonfile.Close()

	// Read data
	b, err := ioutil.ReadAll(jsonfile)
	if err != nil {
		fmt.Println(err)
	}
	data := Data{}
	json.Unmarshal(b, &data)

	// API
	r.GET("/list", func(c *gin.Context) {
		rows := data.ToCalendar()
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
