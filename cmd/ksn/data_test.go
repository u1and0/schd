package ksn

import (
	"encoding/json"
	"testing"
)

func TestData_ToCalendar(t *testing.T) {
	b := ReadJSON("sample.json")
	data := Data{}
	json.Unmarshal(b, &data)
	actual := data.ToCalendar()

	b = ReadJSON("calendar.json")
	expected := Cal{}
	json.Unmarshal(b, &expected)

	for date, idt := range actual {
		for i, actual := range idt.Konpo {
			if expected[date][i] != actual {
				t.Fatalf("got: %v want: %v", actual, expected)
			}
		}
	}
}
