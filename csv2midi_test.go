package main

import "testing"

func TestMakeRandomOffsets(t *testing.T) {
	offsets := makeRandomOffsets(5, 10)
	if len(offsets) != 10 {
		t.Fatalf("expected: %v actual: %v", 10, len(offsets))
	}
	sum := 0
	for _, v := range offsets {
		sum += v
	}

	if sum != 0 {
		t.Fatalf("expected: %v actual: %v", 0, sum)
	}
}
