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

	Cal   map[time.Time]IDmap
	IDmap struct {
		Konpo IDs `json:"梱包ID"`
		Syuka IDs `json:"出荷ID"`
		Noki  IDs `json:"納期ID"`
	}
	IDs []string

	Rows []Row
	Row  struct {
		Date    time.Time `json:"日付"`
		KonpoID string    `json:"梱包-生産番号"`
		// KonpoName   string    `json:"梱包-機器名"`
		// KonpoAssign string    `json:"梱包-担当者"`
		// KonpoIrai   bool      `json:"梱包会社依頼要否"`
		// WDH         string    `json:"外寸法"`
		// Mass        int       `json:"質量"`
		// Yuso        string    `json:"輸送手段"`
		// Chaku       string    `json:"到着予定日"`
		// ToiawaseNo  string    `json:"問合わせ番号"`
		// KonpoMisc   string    `json:"梱包-備考"`

		SyukaID string `json:"出荷-生産番号"`
		// SyukaName   string `json:"出荷-機器名"`
		// SyukaAssign string `json:"出荷-担当者"`
		// SyukaMisc   string `json:"出荷-備考"`

		NokiID string `json:"納期-生産番号"`
		// NokiName   string `json:"納期-機器名"`
		// NokiAssign string `json:"納期-担当者"`
		// NokiMisc   string `json:"納期-備考"`
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
			idmap IDmap
		)

		for id, datum := range *d {
			if date.Equal(datum.Konpo.Date) {
				idmap.Konpo = append(idmap.Konpo, id)
			}
			if date.Equal(datum.Syuka.Date) {
				idmap.Syuka = append(idmap.Syuka, id)
			}
			if date.Equal(datum.Noki.Date) {
				idmap.Noki = append(idmap.Noki, id)
			}
		}
		cal[date] = idmap
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
		col := max(len(ids.Konpo), len(ids.Syuka), len(ids.Noki))
		for i := 0; i <= col; i++ {
			r := Row{Date: date}
			switch {
			case len(ids.Konpo) > i:
				r.KonpoID = ids.Konpo[i]
			case len(ids.Syuka) > i:
				r.SyukaID = ids.Syuka[i]
			case len(ids.Noki) > i:
				r.NokiID = ids.Noki[i]
			}
			rows = append(rows, r)
		}
	}
	return
}

func readJSON(f string) []byte {
	// Open file
	jsonfile, err := os.Open(f)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonfile.Close()

	// Read data
	b, err := ioutil.ReadAll(jsonfile)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func main() {
	r := gin.Default()

	b := readJSON("sample.json")
	data := Data{}
	json.Unmarshal(b, &data)

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
		id := c.Param("id")
		c.JSON(http.StatusOK, data[id])
	})
	r.GET("/all", func(c *gin.Context) {
		c.JSON(http.StatusOK, data)
	})

	r.Run(PORT)
}
