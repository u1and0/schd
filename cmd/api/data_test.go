package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
)

var testdata = ctrl.Data{}

func init() {
	if err := testdata.ReadJSON("../../test/sample.json"); err != nil {
		panic(err)
	}
}

func TestIndex(t *testing.T) {
	expected := testdata
	ts := httptest.NewServer(func() *gin.Engine {
		r := gin.Default()
		r.GET("/api/v1/data/all", Index)
		return r
	}())
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code  200, got %v", resp.StatusCode)
	}
	actual, _ := ioutil.ReadAll(resp.Body)
	if reflect.DeepEqual(actual, expected) {
		t.Fatalf("got: %v\nwant: %v", actual, expected)
	}
}
