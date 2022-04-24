package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/u1and0/schd/cmd/ctrl"
)

const (
	TESTFILE        = "../../test/sample.json"
	TESTADDRESSFILE = "../../test/address.json"
	URL             = "http://localhost:8080/api/v1/data"
	ID              = "741744"
)

var testdata ctrl.Data

func init() {
	if err := testdata.ReadJSON(TESTFILE); err != nil {
		panic(err)
	}
}

func rollbackTestfile() {
	if _, err := exec.Command("cp", "-f", "../../test/sample_fix.json", TESTFILE).Output(); err != nil {
		panic(err)
	}
}

/*すべてのtestはgo run main.go して、
アプリケーションを立ち上げた状態でtestする
httptest 使ったが、よくわからなかった...
*/

func TestIndex(t *testing.T) {
	expected := testdata
	resp, err := http.Get(URL + "/")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code  200, got %v", resp.StatusCode)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	actual := ctrl.Data{}
	json.Unmarshal(b, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got: %#v\nwant: %#v", actual, expected)
	}
}

func TestGet(t *testing.T) {
	expected := testdata[ctrl.ID(ID)]
	resp, err := http.Get(URL + "/" + ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code  200, got %v", resp.StatusCode)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	// id 指定なので、返ってくるのはDataではなく、Datum
	actual := ctrl.Datum{}
	json.Unmarshal(b, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got: %#v\nwant: %#v", actual, expected)
	}
}

func TestPost(t *testing.T) {
	konpo := ctrl.Konpo{
		Date:       time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC),
		Irai:       "否",
		WDH:        "2100x2190x1560",
		Mass:       32,
		Yuso:       "宅急便",
		Chaku:      time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC),
		ToiawaseNo: "123456",
		Misc:       "備考欄",
	}
	syuka := ctrl.Syuka{
		Date: time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC),
		Misc: "",
	}
	noki := ctrl.Noki{
		Date: time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC),
		Misc: "",
	}
	datum := ctrl.Datum{
		Name:   "B-GHT-222",
		Assign: "Putin",
		Konpo:  konpo,
		Syuka:  syuka,
		Noki:   noki,
	}
	expected := ctrl.Data{"990001": datum}

	s := `{
    "990001":{
        "機器名": "B-GHT-222",
        "担当者": "Putin",
        "梱包": {
            "日付": "2022-04-01T00:00:00Z",
            "梱包会社依頼要否": "否",
            "外寸法": "2100x2190x1560",
            "質量": 32,
            "輸送手段": "宅急便",
            "到着予定日": "2022-04-04T00:00:00Z",
            "問合わせ番号": "123456",
            "備考": "備考欄"
        },
        "出荷日": {
            "日付": "2022-04-14T00:00:00Z",
            "備考": ""
        },
        "納期": {
            "日付": "2022-04-14T00:00:00Z",
            "備考": ""
        }
    }
}`
	resBody := strings.NewReader(s)
	resp, err := http.Post(URL+"/add", "application/json", resBody)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code  200, got %v", resp.StatusCode)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	actual := ctrl.Data{}
	json.Unmarshal(b, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Log(string(b))
		t.Errorf("got: %#v\nwant: %#v", actual, expected)
	}
	rollbackTestfile()
}

func TestDelete(t *testing.T) {
	expected := fmt.Sprintf(`{"id": %v"}`, ID)
	resBody := strings.NewReader(expected)
	resp, err := http.Post(URL+"/"+ID, "application/json", resBody)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code  204, got %v", resp.StatusCode)
	}
	actual, _ := ioutil.ReadAll(resp.Body)
	if string(actual) != expected {
		t.Errorf("got: %#v\nwant: %#v", string(actual), expected)
	}
	rollbackTestfile()
}

func TestFetchAddress(t *testing.T) {
	// Read JSON
	var expected AddressMap
	var actual AddressMap
	f, err := os.Open(TESTADDRESSFILE)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	err = json.Unmarshal(b, &expected)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Fetch JSON
	resp, err := http.Get(URL + "/address")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code  200, got %v", resp.StatusCode)
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	json.Unmarshal(r, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got: %#v\nwant: %#v", actual, expected)
	}
}
