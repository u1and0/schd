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
	s.Width = []int{10, 20, 30, 40}
	s.Length = []int{100, 200, 300, 400}
	s.Hight = []int{1000, 2000, 3000, 4000}
	s.Mass = []int{1, 2, 3, 4}
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

func TestSize_String(t *testing.T) {
	actual := s.ToString()
	expected := `10x100x1000mm 1kg, 20x200x2000mm 2kg, 30x300x3000mm 3kg, 40x400x4000mm 4kg`
	if actual != expected {
		t.Fatalf("got: %v want: %v", actual, expected)
	}
}
