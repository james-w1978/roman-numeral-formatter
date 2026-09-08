package main

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		value   int
		wantErr bool
	}{
		{"xiv", "XIV", 14, false},
		{"  MCMxciv  ", "MCMXCIV", 1994, false},
		{"IIII", "IV", 4, false},
		{"", "", 0, true},
		{"XIIB", "", 0, true},
		{"MMMM", "", 0, true}, // 4000, out of range
	}

	for _, c := range cases {
		got, value, err := Normalize(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("Normalize(%q): expected error, got %q", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("Normalize(%q): unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want || value != c.value {
			t.Errorf("Normalize(%q) = (%q, %d), want (%q, %d)", c.in, got, value, c.want, c.value)
		}
	}
}
