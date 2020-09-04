package main

import (
	"testing"
)

func TestSelectColorIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int
		out  int
	}{
		{
			"zero",
			0,
			1,
		},
		{
			"one",
			1,
			2,
		},
		{
			"two",
			2,
			3,
		},
		{
			"three",
			3,
			0,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			out := selectColorIndex(v.in)
			if out != v.out {
				t.Errorf("input %v\n, get: %#v\n, want: %#v\n", v.in, out, v.out)
			}
		})
	}
}
