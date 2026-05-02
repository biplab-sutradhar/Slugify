package idgen

import (
	"testing"
)

func TestEncode_Zero(t *testing.T) {
	got := Encode(0)
	if got == "" {
		t.Fatal("Encode(0) returned empty string")
	}
	if len(got) != 1 {
		t.Errorf("Encode(0) length = %d, want 1", len(got))
	}
}

func TestEncode_Deterministic(t *testing.T) {
	a := Encode(123456)
	b := Encode(123456)
	if a != b {
		t.Errorf("Encode is non-deterministic: %q vs %q", a, b)
	}
}

func TestEncode_DifferentInputsProduceDifferentOutputs(t *testing.T) {
	tests := []int64{1, 2, 61, 62, 63, 1000, 100_000, 1_000_000_000}
	seen := make(map[string]int64)
	for _, n := range tests {
		got := Encode(n)
		if prev, dup := seen[got]; dup {
			t.Errorf("Encode collision: %d and %d both => %q", prev, n, got)
		}
		seen[got] = n
	}
}

func TestEncode_LengthGrowsWithMagnitude(t *testing.T) {
	cases := []struct {
		id     int64
		minLen int
		maxLen int
	}{
		{1, 1, 1},
		{61, 1, 1},
		{62, 2, 2},
		{62 * 62, 3, 3},
		{1_000_000_000, 5, 6}, // ~5–6 chars
	}
	for _, c := range cases {
		got := Encode(c.id)
		if len(got) < c.minLen || len(got) > c.maxLen {
			t.Errorf("Encode(%d) = %q (len %d); want length in [%d, %d]",
				c.id, got, len(got), c.minLen, c.maxLen)
		}
	}
}
