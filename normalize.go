package main

import (
	"fmt"
	"strings"
)

// symbolValues maps each Roman numeral letter to its integer value.
var symbolValues = map[byte]int{
	'I': 1,
	'V': 5,
	'X': 10,
	'L': 50,
	'C': 100,
	'D': 500,
	'M': 1000,
}

// canonicalTable lists numerals in descending order, including the standard
// subtractive pairs, so encode can build the canonical string greedily.
var canonicalTable = []struct {
	value  int
	symbol string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
	{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
	{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

const (
	minValue = 1
	maxValue = 3999 // no standard single-symbol notation above this
)

// Normalize takes messy Roman numeral text (stray whitespace, mixed case,
// non-canonical repetition like IIII) and returns the canonical spelling
// plus the integer value it represents.
func Normalize(input string) (string, int, error) {
	return normalize(input, false)
}

// NormalizeStrict behaves like Normalize but rejects any input that isn't
// already in canonical form, rather than silently correcting it.
func NormalizeStrict(input string) (string, int, error) {
	return normalize(input, true)
}

func normalize(input string, strict bool) (string, int, error) {
	cleaned := strings.ToUpper(strings.Join(strings.Fields(input), ""))
	if cleaned == "" {
		return "", 0, fmt.Errorf("empty input")
	}

	value := 0
	for i := 0; i < len(cleaned); i++ {
		v, ok := symbolValues[cleaned[i]]
		if !ok {
			return "", 0, fmt.Errorf("invalid character %q in %q", cleaned[i], input)
		}
		// A smaller value immediately before a larger one is subtracted
		// (IV = 4); everything else is added. Treating every other pair as
		// additive is what lets us tolerate sloppy input like IIII or VIIII
		// instead of rejecting anything that isn't already canonical.
		if i+1 < len(cleaned) {
			if next, ok := symbolValues[cleaned[i+1]]; ok && v < next {
				value -= v
				continue
			}
		}
		value += v
	}

	if value < minValue || value > maxValue {
		return "", 0, fmt.Errorf("value %d out of range %d-%d", value, minValue, maxValue)
	}

	canonical := encode(value)
	if strict && canonical != cleaned {
		return "", 0, fmt.Errorf("%q is not canonical (expected %q)", input, canonical)
	}

	return canonical, value, nil
}

// ToRoman converts an integer into its canonical Roman numeral spelling.
func ToRoman(value int) (string, error) {
	if value < minValue || value > maxValue {
		return "", fmt.Errorf("value %d out of range %d-%d", value, minValue, maxValue)
	}
	return encode(value), nil
}

// encode converts an integer into canonical Roman numeral form.
func encode(value int) string {
	var b strings.Builder
	for _, entry := range canonicalTable {
		for value >= entry.value {
			b.WriteString(entry.symbol)
			value -= entry.value
		}
	}
	return b.String()
}
