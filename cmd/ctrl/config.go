package ctrl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// CONFIGPATH : 設定ファイルパス
const CONFIGPATH = "db/config.json"

// Config : 設定json
var Config ConfigType

type (
	// ConfigType : configuration
	ConfigType struct {
		AllocatePath string `json:"配車要求票パス"`
		// AddressMap : JSON ファイルから読み取った住所録
		AddressMap map[string]Address `json:"住所録"`
		// Section 所属課のmap
		Section map[string]string `json:"配車要求票頭番号"`
	}
	// Address : 1住所あたり5行まで, 1行あたり15文字まで
	// Excelシートの枠の都合
	Address []string
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

func (a *Address) String() string {
	return strings.Join(*a, "\n")
}
