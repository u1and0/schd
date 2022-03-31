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
	if err := expected.ReadJSON("../../test/calendar.json"); err != nil {
		t.Fatal(err)
	}
	if !(reflect.DeepEqual(actual, expected)) {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

func TestData_Unstack(t *testing.T) {
	actual := testdata.Stack().Unstack()
	expected := Rows{}
	if err := expected.ReadJSON("../../test/rows.json"); err != nil {
		t.Fatal(err)
	}
	for i, v := range expected {
		if actual[i] != v {
			t.Errorf("got: %v\nwant: %v", actual, expected)
		}
	}
}
