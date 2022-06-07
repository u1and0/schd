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

type (
	// ConfigType : configuration
	ConfigType struct {
		AllocatePath string `json:"配車要求票パス"`
		// Section 所属課のmap
		Section map[string]string `json:"配車要求票頭番号"`
	}
)

func init() {
	if err := UnmarshalJSONfile(&Config, CONFIGPATH); err != nil {
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

// UnmarshalJSONfile : Read json from .json file then unmarshal as some T type
func UnmarshalJSONfile(T interface{}, filename string) error {
	b, err := readFile(filename)
	// As JSON
	err = json.Unmarshal(b, &T)
	return err
}
