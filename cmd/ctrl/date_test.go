package ctrl

import (
	"testing"
	"time"
)

func Test_ToMonth(t *testing.T) {
	actualy, actualm, err := ToMonth("2204")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedy := 2022
	if actualy != expectedy {
		t.Errorf("got: %v want: %v", actualy, expectedy)
	}
	expectedm := 4
	if actualm != expectedm {
		t.Errorf("got: %v want: %v", actualm, expectedm)
	}
}

func Test_DayofFirstEnd(t *testing.T) {
	actualf, actuale := DayofFirstEnd(2022, 4)
	expectedf, expectede := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 4, 30, 0, 0, 0, 0, time.UTC)
	if actualf != expectedf {
		t.Errorf("got: %v want: %v", actualf, expectedf)
	}
	if actuale != expectede {
		t.Errorf("got: %v want: %v", actuale, expectede)
	}
	actualf, actuale = DayofFirstEnd(2020, 2) // うるう年
	expectedf, expectede = time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC)
	if actualf != expectedf {
		t.Errorf("got: %v want: %v", actualf, expectedf)
	}
	if actuale != expectede {
		t.Errorf("got: %v want: %v", actuale, expectede)
	}
}

// func Test_Iter(t *testing.T) {
// 	st := time.Date(2022, 4, 1)
// 	st := time.Date(2022, 4, 1)
// 	actual := Iter()
// 	expected :=
// 		t.Fatal(err)
// }
