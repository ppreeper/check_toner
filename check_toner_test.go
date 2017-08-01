package main

import (
	"testing"
)

// mulColorTests use values that you know are right
var mulColorTests = []struct {
	color    string
	maxValue string
	lvlValue string
	expected string
}{
	{"1d6", "d", "", ""},
}

// TestTonerOutput test
func TestTonerOutput(t *testing.T) {
	for _, mt := range mulColorTests {
		tonerOutput("K", "100", "90")
		t.Errorf("%v", mt.color)
	}
}

// BenchmarkPattern
func BenchmarkTonerOutput(b *testing.B) {
	// run the Fib function b.N times
}
