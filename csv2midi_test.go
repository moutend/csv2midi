package main

import "testing"

func TestRandomizer_Randomize(t *testing.T) {
	r := &randomizer{
		factor:   10,
		position: 0,
	}
	sum := 0
	for i := 0; i < 100; i++ {
		v := r.Randomize(i)
		if v < 0 {
			t.Fatal("expected: 0 actual: %v", v)
		}
		d := v - i
		sum += d
	}
	if sum-r.position != 0 {
		t.Fatalf("expected: %v actual: %v", 0, sum-r.position)
	}
}
