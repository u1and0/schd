package ctrl

import (
	"errors"
	"fmt"
	"time"
)

type (
	// Datum : 生産番号ごとの物流情報
	Datum struct {
		Name   string `json:"機器名" form:"name"`
		Assign string `json:"担当者" form:"assign"`
		Konpo  `json:"梱包"`
		Syuka  `json:"出荷日"`
		Noki   `json:"納期"`
	}
	// Konpo : 梱包列情報
	Konpo struct {
		Date       time.Time `json:"日付" form:"konpo-date" time_format:"2006/01/02"`
		Irai       string    `json:"梱包会社依頼要否" form:"irai"`
		WDH        string    `json:"外寸法" form:"wdh"`
		Mass       int       `json:"質量" form:"mass"`
		Yuso       string    `json:"輸送手段" form:"yuso"`
		Chaku      time.Time `json:"到着予定日" form:"chaku" time_format:"2006/01/02"`
		ToiawaseNo string    `json:"問合わせ番号" form:"toiawase-no"`
		Misc       string    `json:"備考" form:"konpo-misc"`
	}
	// Syuka : 出荷列情報
	Syuka struct {
		Date time.Time `json:"日付" form:"syuka-date" time_format:"2006/01/02"`
		Misc string    `json:"備考" form:"syuka-misc"`
	}
	// Noki : 納期列情報
	Noki struct {
		Date time.Time `json:"日付" form:"noki-date" time_format:"2006/01/02"`
		Misc string    `json:"備考" form:"noki-misc"`
	}
)

func (up *Datum) Update(id ID, data *Data) error {
	// Data exist check
	if _, ok := (*data)[id]; !ok {
		msg := fmt.Sprintf("ID: %v データが存在しません。/api/v1/data/addを試してください。", id)
		return errors.New(msg)
	}
	// Update data
	(*data)[id] = *up
	return nil
}
