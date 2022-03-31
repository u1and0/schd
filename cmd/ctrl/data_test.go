package ctrl

import (
	"reflect"
	"testing"
)

var testdata = Data{}

func init() {
	if err := testdata.ReadJSON("../../test/sample.json"); err != nil {
		panic(err)
	}
}

func TestData_Stack(t *testing.T) {
	actual := testdata.Stack()
	expected := Cal{}
	err := expected.ReadJSON("../../test/calendar.json")
	if err != nil {
		panic(err)
	}
	if !(reflect.DeepEqual(actual, expected)) {
		t.Fatalf("got: %v\nwant: %v", actual, expected)
	}
}
