package ctrl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type (
	// Data : 生産番号により分けられた物流情報
	Data map[ID]Datum
	// Marshaler : JSON reader/writer
	Marshaler interface {
		ReadJSON(fs string)
		WriteJSON(fs string)
	}
)

// ReadJSON : Read from json file to Data structure
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

// WriteJSON : Write to json file from Data structure
func (d *Data) WriteJSON(fs string) error {
	// To binary
	b, err := json.MarshalIndent(&d, "", "\t")
	if err != nil {
		return err
	}
	// Write file
	err = ioutil.WriteFile(fs, b, 0644)
	return err
}

func (ad *Data) Add(data *Data) error {
	// Data exist check
	for k := range *ad {
		if _, ok := (*data)[k]; ok {
			msg := fmt.Sprintf("ID: %v データが既に存在しています。Updateを試してください。", k)
			return errors.New(msg)
		}
	}
	for k, v := range *ad {
		(*data)[k] = v
	}
	return nil
}

func (id ID) Del(data *Data) error {
	if _, ok := (*data)[id]; !ok {
		msg := fmt.Sprintf("ID: %v が見つかりません。別のIDを指定してください。", id)
		return errors.New(msg)
	}
	delete((*data), id)
	// Check deleted id
	if _, ok := (*data)[id]; ok {
		msg := fmt.Sprintf("ID: %v を削除できませんでした。", id)
		return errors.New(msg)
	}
	return nil
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
