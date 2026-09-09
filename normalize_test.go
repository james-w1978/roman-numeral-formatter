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

func TestToRoman(t *testing.T) {
	cases := []struct {
		in      int
		want    string
		wantErr bool
	}{
		{14, "XIV", false},
		{1994, "MCMXCIV", false},
		{4, "IV", false},
		{1, "I", false},
		{3999, "MMMCMXCIX", false},
		{0, "", true},
		{4000, "", true},
		{-5, "", true},
	}

	for _, c := range cases {
		got, err := ToRoman(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("ToRoman(%d): expected error, got %q", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ToRoman(%d): unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ToRoman(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestToRomanNormalizeRoundTrip(t *testing.T) {
	for value := minValue; value <= maxValue; value++ {
		roman, err := ToRoman(value)
		if err != nil {
			t.Fatalf("ToRoman(%d): unexpected error: %v", value, err)
		}
		gotRoman, gotValue, err := Normalize(roman)
		if err != nil {
			t.Fatalf("Normalize(%q): unexpected error: %v", roman, err)
		}
		if gotRoman != roman || gotValue != value {
			t.Errorf("round trip for %d: got (%q, %d), want (%q, %d)", value, gotRoman, gotValue, roman, value)
		}
	}
}
