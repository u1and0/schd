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

	Cal map[time.Time]IDs
	IDs struct {
		Konpo []string `json:"梱包ID"`
		Syuka []string `json:"出荷ID"`
		Noki  []string `json:"納期ID"`
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
func (d *Data) ToCalendar() Cal {
	cal := Cal{}
	dates := []time.Time{
		time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 24, 0, 0, 0, 0, time.UTC),
	}
	for _, date := range dates {
		var (
			ids IDs
		)

		for id, datum := range *d {
			if date.Equal(datum.Konpo.Date) {
				ids.Konpo = append(ids.Konpo, id)
			}
			if date.Equal(datum.Syuka.Date) {
				ids.Syuka = append(ids.Syuka, id)
			}
			if date.Equal(datum.Noki.Date) {
				ids.Noki = append(ids.Noki, id)
			}
		}
		cal[date] = ids
	}
	return cal
}

func max(s ...int) (x int) {
	for _, i := range s {
		if x < i {
			x = i
		}
	}
	return
}

func (c Cal) ToRows() (rows Rows) {
	for date, ids := range c {
		r := Row{Date: date}
		col := max(len(ids.Konpo), len(ids.Syuka), len(ids.Noki))
		for i := 0; i < col; i++ {
			if len(ids.Konpo) > i {
				r.KonpoID = ids.Konpo[i]
			} else {
				r.KonpoID = ""
			}
			if len(ids.Syuka) > i {
				r.SyukaID = ids.Syuka[i]
			} else {
				r.SyukaID = ""
			}
			if len(ids.Noki) > i {
				r.NokiID = ids.Noki[i]
			} else {
				r.NokiID = ""
			}
			rows = append(rows, r)
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
		rows := data.ToCalendar().ToRows()
		c.JSON(http.StatusOK, rows)
	})
	r.GET("/cal", func(c *gin.Context) {
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
