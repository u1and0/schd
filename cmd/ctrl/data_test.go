package ctrl

import (
	"reflect"
	"testing"
	"time"
)

var testdata = Data{}

func init() {
	if err := testdata.ReadJSON("../../test/sample.json"); err != nil {
		panic(err)
	}
}

func TestData_Stack(t *testing.T) {
	actual := testdata.Stack()
	expected := Cal{
		time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC): IDt{
			Konpo: IDs{"741744"},
		},
		time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC): IDt{
			Konpo: IDs{"739111"},
			Syuka: IDs{"741744"},
		},
		time.Date(2022, 4, 20, 0, 0, 0, 0, time.UTC): IDt{},
		time.Date(2022, 4, 24, 0, 0, 0, 0, time.UTC): IDt{
			Syuka: IDs{"739111"},
			Noki:  IDs{"739111", "741744"},
		},
	}
	if !(reflect.DeepEqual(actual, &expected)) {
		t.Errorf("got: %v\nwant: %v", actual, &expected)
	}
}
