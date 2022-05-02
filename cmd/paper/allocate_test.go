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
	}
}

func TestSize_Sum(t *testing.T) {
	actual := s.Sum()
	expected := 13
	if actual != expected {
		t.Fatalf("got: %v want: %v", actual, expected)
	}
}

func TestSize_Compile(t *testing.T) {
	actual := s.Compile()
	expected := PackageCount{"みかん": 8, "りんご": 4, "すいか": 1}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got: %#v want: %#v", actual, expected)
	}
}

func TestPackageCount_String(t *testing.T) {
	actual := p.ToString()
	expected1 := `りんご(4), みかん(8)`
	expected2 := `みかん(8), りんご(4)`
	if actual != expected1 && actual != expected2 {
		t.Fatalf("got: %v want: %v or %v", actual, expected1, expected2)
	}
}
