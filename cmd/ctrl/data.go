package ctrl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type (
	// Data : map of Datum
	Data map[ID]Datum
	// Datum : date info
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

	Cal map[time.Time]IDt
	IDt struct {
		Konpo IDs `json:"梱包ID"`
		Syuka IDs `json:"出荷ID"`
		Noki  IDs `json:"納期ID"`
	}
	IDs []ID
	ID  string

	Rows []Row
	Row  struct {
		Date    time.Time `json:"日付"`
		KonpoID ID        `json:"梱包-生産番号"`
		// KonpoName   string    `json:"梱包-機器名"`
		// KonpoAssign string    `json:"梱包-担当者"`
		// KonpoIrai   bool      `json:"梱包会社依頼要否"`
		// WDH         string    `json:"外寸法"`
		// Mass        int       `json:"質量"`
		// Yuso        string    `json:"輸送手段"`
		// Chaku       string    `json:"到着予定日"`
		// ToiawaseNo  string    `json:"問合わせ番号"`
		// KonpoMisc   string    `json:"梱包-備考"`

		SyukaID ID `json:"出荷-生産番号"`
		// SyukaName   string `json:"出荷-機器名"`
		// SyukaAssign string `json:"出荷-担当者"`
		// SyukaMisc   string `json:"出荷-備考"`

		NokiID ID `json:"納期-生産番号"`
		// NokiName   string `json:"納期-機器名"`
		// NokiAssign string `json:"納期-担当者"`
		// NokiMisc   string `json:"納期-備考"`
	}
)

func (d *Data) ReadJSON(fs string) error {
	// Open file
	f, err := os.Open(fs)
	defer f.Close()
	if err != nil {
		return err
	}

	// Read data
	b, err := ioutil.ReadAll(f)
	json.Unmarshal(b, &d)
	return err
}

// Stack : 製番jsonを走査し、
// 日付をキーに、項目ごとに製番リストを保持する
// Cal構造体を返す
func (d *Data) Stack() Cal {
	cal := Cal{}
	dates := []time.Time{
		time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 4, 24, 0, 0, 0, 0, time.UTC),
	}
	for _, date := range dates {
		var idt IDt
		for id, datum := range *d {
			if date.Equal(datum.Konpo.Date) {
				idt.Konpo = append(idt.Konpo, id)
			}
			if date.Equal(datum.Syuka.Date) {
				idt.Syuka = append(idt.Syuka, id)
			}
			if date.Equal(datum.Noki.Date) {
				idt.Noki = append(idt.Noki, id)
			}
		}
		cal[date] = idt
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

// Unstack : Cal構造体から
// 日付をプライマリキーとするテーブル形式のRowsを返す
func (c Cal) Unstack() (rows Rows) {
	for date, idt := range c {
		l := max(len(idt.Konpo), len(idt.Syuka), len(idt.Noki))
		// 何もない日でも一行は空行出力
		if l == 0 {
			r := Row{Date: date}
			rows = append(rows, r)
			continue
		}
		for i := 0; i < l; i++ {
			r := Row{Date: date}
			if len(idt.Konpo) > i {
				r.KonpoID = idt.Konpo[i]
			} else {
				r.KonpoID = ""
			}
			if len(idt.Syuka) > i {
				r.SyukaID = idt.Syuka[i]
			} else {
				r.SyukaID = ""
			}
			if len(idt.Noki) > i {
				r.NokiID = idt.Noki[i]
			} else {
				r.NokiID = ""
			}
			rows = append(rows, r)
		}
	}
	return
}
