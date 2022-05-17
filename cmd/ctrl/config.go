package ctrl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// CONFIGPATH : 設定ファイルパス
const CONFIGPATH = "db/config.json"

// Config : 設定json
var Config ConfigType

// ConfigType : configuration
type ConfigType struct {
	AllocatePath string `json:"配車要求票パス"`
}

func init() {
	if err := UnmarshalJSON(&Config, CONFIGPATH); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", Config)
}

// readFile : read filepath to binary
func readFile(fs string) (b []byte, err error) {
	// Open file
	f, err := os.Open(fs)
	defer f.Close()
	if err != nil {
		return
	}
	// Read file
	b, err = ioutil.ReadAll(f)
	if err != nil {
		return
	}
	return
}

// UnmarshalJSON : some T type
func UnmarshalJSON(T interface{}, filename string) error {
	b, err := readFile(filename)
	// As JSON
	err = json.Unmarshal(b, &T)
	return err
}
