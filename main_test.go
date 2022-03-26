package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

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

func TestData_ToCalendar(t *testing.T) {
	b := readJSON("sample.json")
	data := Data{}
	json.Unmarshal(b, &data)
	actual := data.ToCalendar()

	b = readJSON("ToCalendar.json")
	expected := Rows{}
	json.Unmarshal(b, &expected)

	for i, obj := range actual {
		if expected[i] != obj {
			t.Fatalf("got: %v want: %v", actual, expected)
		}
	}
}
