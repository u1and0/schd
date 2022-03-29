package ksn

import (
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

func (ids *IDs) Append(id ID) {
	*ids = append(*ids, id)
}

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
		var idt IDt
		for id, datum := range *d {
			switch {
			case datum.Konpo.Date.Equal(date):
				idt.Konpo = append(idt.Konpo, id)
				fallthrough
			case datum.Syuka.Date.Equal(date):
				idt.Syuka = append(idt.Syuka, id)
				fallthrough
			case datum.Noki.Date.Equal(date):
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
