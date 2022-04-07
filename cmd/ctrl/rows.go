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
		Date        time.Time `json:"日付"`
		KonpoID     ID        `json:"梱包-生産番号"`
		KonpoName   string    `json:"梱包-機器名"`
		KonpoAssign string    `json:"梱包-担当者"`
		Irai        string    `json:"梱包会社依頼要否"`
		WDH         string    `json:"外寸法"`
		Mass        int       `json:"質量"`
		Yuso        string    `json:"輸送手段"`
		Chaku       time.Time `json:"到着予定日" time_format:"2006/01/02"`
		ToiawaseNo  string    `json:"問合わせ番号"`
		KonpoMisc   string    `json:"梱包-備考"`

		SyukaID     ID     `json:"出荷-生産番号"`
		SyukaName   string `json:"出荷-機器名"`
		SyukaAssign string `json:"出荷-担当者"`
		SyukaMisc   string `json:"出荷-備考"`

		NokiID     ID     `json:"納期-生産番号"`
		NokiName   string `json:"納期-機器名"`
		NokiAssign string `json:"納期-担当者"`
		NokiMisc   string `json:"納期-備考"`
	}
)

// ReadJSON : Read from json file to Rows structure
func (r *Rows) ReadJSON(fs string) error {
	// Open file
	f, err := os.Open(fs)
	defer f.Close()
	if err != nil {
		return err
	}
	// Read data
	b, err := ioutil.ReadAll(f)
	json.Unmarshal(b, &r)
	return err
}

// Verbose : Rows の欠落情報をdata[id]から補完
func (r *Rows) Verbose(d Data) Rows {
	v := *r
	for i, row := range *r {
		id := row.KonpoID
		v[i].KonpoName = d[id].Name
		v[i].KonpoAssign = d[id].Assign
		v[i].Irai = d[id].Irai
		v[i].WDH = d[id].WDH
		v[i].Mass = d[id].Mass
		v[i].Yuso = d[id].Yuso
		v[i].Chaku = d[id].Chaku
		v[i].ToiawaseNo = d[id].ToiawaseNo
		v[i].KonpoMisc = d[id].Konpo.Misc

		id = row.SyukaID
		v[i].SyukaName = d[id].Name
		v[i].SyukaAssign = d[id].Assign
		v[i].SyukaMisc = d[id].Syuka.Misc

		id = row.NokiID
		v[i].NokiName = d[id].Name
		v[i].NokiAssign = d[id].Assign
		v[i].NokiMisc = d[id].Noki.Misc
	}
	return v
}
