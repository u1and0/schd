package main

import (
	"encoding/json"
	"testing"
)

func TestData_ToCalendar(t *testing.T) {
	b := readJSON("sample.json")
	data := Data{}
	json.Unmarshal(b, &data)
	actual := data.ToCalendar()

	b = readJSON("calendar.json")
	expected := Rows{}
	json.Unmarshal(b, &expected)

	for i, obj := range actual {
		if expected[i] != obj {
			t.Fatalf("got: %v want: %v", actual, expected)
		}
	}
}
