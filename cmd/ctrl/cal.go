package ctrl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type (
	// Cal : Idtを日付ごとにまとめたmap
	Cal map[time.Time]IDt
	// IDt : 列情報
	IDt struct {
		Konpo IDs
		Syuka IDs
		Noki  IDs
	}
	// IDs : ID のスライス
	IDs []ID
	// ID : 生産番号
	// 数字6桁, ただしstring型
	// JSON のキーがstringしか受け付けないため。
	ID string
)

// ReadJSON : Read from json file to Cal structure
func (d *Cal) ReadJSON(fs string) error {
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

// Unstack : Cal構造体から
// 日付をプライマリキーとするテーブル形式のRowsを返す
func (d Cal) Unstack() (rows Rows) {
	for date, idt := range d {
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

func max(s ...int) (x int) {
	for _, i := range s {
		if x < i {
			x = i
		}
	}
	return
}
