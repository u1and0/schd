package paper

import (
	"reflect"
	"testing"
)

var (
	s Size
	p PackageCount
)

func init() {
	s.Package = []string{"みかん", "りんご", "みかん", "すいか"}
	s.Quantity = []int{5, 4, 3, 1}
	p = PackageCount{
		"りんご": 4,
		"みかん": 8,
		"すいか": 1,
	}
}

func TestSize_Sum(t *testing.T) {
}

func TestSize_Compile(t *testing.T) {
	actual := s.Compile()
	expected := p
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got: %#vwant: %#v", actual, expected)
	}
}

func TestPackageCount_String(t *testing.T) {
	actual := p.ToString()
	expected := `りんご(4)
みかん(8)
すいか(1)`
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got: %#v want: %#v", actual, expected)
	}
}
