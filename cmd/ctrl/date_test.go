package ctrl

import (
	"testing"
)

func Test_ToMonth(t *testing.T) {
	actualy, actualm, err := ToMonth("2204")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedy := 2022
	if actualy != expectedy {
		t.Errorf("got: %vwant: %v", actualy, expectedy)
	}
	expectedm := 4
	if actualm != expectedm {
		t.Errorf("got: %vwant: %v", actualm, expectedm)
	}
}

// func Test_Iter(t *testing.T) {
// 	st := time.Date(2022, 4, 1)
// 	st := time.Date(2022, 4, 1)
// 	actual := Iter()
// 	expected :=
// 		t.Fatal(err)
// }
