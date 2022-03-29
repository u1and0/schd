package ksn

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadJSON(fs string) []byte {
	// Open file
	f, err := os.Open(fs)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	// Read data
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}
	return b
}
