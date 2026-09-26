package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadStdinRecords(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		delimiter string
		want      []string
	}{
		{"newline default", "xiv\nIIII\nxl\n", "\n", []string{"xiv", "IIII", "xl"}},
		{"comma", "xiv,IIII,xl", ",", []string{"xiv", "IIII", "xl"}},
		{"escaped tab", "xiv\tIIII\txl", `\t`, []string{"xiv", "IIII", "xl"}},
		{"blank records dropped", "xiv,,xl,", ",", []string{"xiv", "xl"}},
		{"multi-char delimiter", "xiv::IIII::xl", "::", []string{"xiv", "IIII", "xl"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := readStdinRecords(strings.NewReader(c.input), c.delimiter)
			if err != nil {
				t.Fatalf("readStdinRecords(%q, %q): unexpected error: %v", c.input, c.delimiter, err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("readStdinRecords(%q, %q) = %v, want %v", c.input, c.delimiter, got, c.want)
			}
		})
	}
}
