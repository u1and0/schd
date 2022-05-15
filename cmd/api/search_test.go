package api

import (
	"testing"
)

func TestSize_ContainsAll(t *testing.T) {
	s := `この文書に含まれている単語すべてが一致したらtrueを返す`
	// truthy test
	keywd := []string{"文書", "すべて", "単語"}
	actual := ContainsAll(s, keywd...)
	expected := true
	if !actual {
		t.Fatalf("got: %v want: %v", actual, expected)
	}
	// falsy test
	keywd = []string{"文書", "すべて", "単語", "一致しない"}
	actual = ContainsAll(s, keywd...)
	expected = false
	if actual {
		t.Fatalf("got: %v want: %v", actual, expected)
	}
}
