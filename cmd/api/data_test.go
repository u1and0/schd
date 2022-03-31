package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/u1and0/schd/cmd/ctrl"
)

var testdata ctrl.Data

func init() {
	if err := testdata.ReadJSON("../../test/sample.json"); err != nil {
		panic(err)
	}
}

/*すべてのtestはgo run main.go して、
アプリケーションを立ち上げた状態でtestする
httptest 使ったが、よくわからなかった...
*/

func TestIndex(t *testing.T) {
	expected := testdata
	resp, err := http.Get("http://localhost:8080/api/v1/data/all")
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
	id := ctrl.ID("741744")
	expected := testdata[id]
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/v1/data/%s", id))
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
		t.Fatalf("got: %#v\nwant: %#v", actual, expected)
	}
}
