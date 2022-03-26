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
		var (
			i   int  // Added row
			a   bool // Pass fallthrough
			row Row
		)
		for id, datum := range *d {
			fmt.Println("date: ", date)
			if date.Equal(datum.Konpo.Date) {
				fmt.Printf("Konpo date 追加: %s,%s\n", id, datum.Konpo.Date)
				fmt.Println("date: ", date)
				row.KonpoID = id
				row.KonpoName = datum.Name
				row.KonpoAssign = datum.Assign
				row.WDH = datum.WDH
				a = true
			}
			if date.Equal(datum.Syuka.Date) {
				fmt.Printf("Syuka date 追加: %s,%s\n", id, datum.Syuka.Date)
				fmt.Println("date: ", date)
				row.SyukaID = id
				row.SyukaName = datum.Name
				row.SyukaAssign = datum.Assign
				a = true
			}
			if date.Equal(datum.Noki.Date) {
				fmt.Printf("Noki date 追加: %s,%s\n", id, datum.Noki.Date)
				fmt.Println("date: ", date)
				row.NokiID = id
				row.NokiName = datum.Name
				row.NokiAssign = datum.Assign
				a = true
			}
			row.Date = date
			// fallthrough を通過してRowに追加したらa==true
			// fallthrough 通過しなくても、date==0であれば、空の行追加
		}
		if a || i < 1 {
			rows = append(rows, row)
			i++ // added row
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
