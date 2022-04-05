package ctrl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type (
	// Rows : HTMLテーブル形式表示用
	Rows []Row
	// Row : Rowsの内の1行
	Row struct {
		Date    time.Time `json:"日付"`
		KonpoID ID        `json:"梱包-生産番号"`
		// KonpoName   string    `json:"梱包-機器名"`
		// KonpoAssign string    `json:"梱包-担当者"`
		// Irai   bool      `json:"梱包会社依頼要否"`
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

// ReadJSON : Read from json file to Rows structure
func (d *Rows) ReadJSON(fs string) error {
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
