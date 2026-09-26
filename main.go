package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type result struct {
	Input      string `json:"input"`
	Normalized string `json:"normalized,omitempty"`
	Value      int    `json:"value,omitempty"`
	Error      string `json:"error,omitempty"`
}

func main() {
	jsonOutput := flag.Bool("json", false, "emit results as JSON, one object per line")
	toRoman := flag.Bool("to-roman", false, "treat input as decimal numbers and convert to Roman numerals")
	strict := flag.Bool("strict", false, "reject non-canonical numerals instead of correcting them")
	delimiter := flag.String("delimiter", "\n", "record separator used when reading stdin")
	flag.Parse()

	lines := flag.Args()
	if len(lines) == 0 {
		if *delimiter == "" {
			fmt.Fprintln(os.Stderr, "romanfmt: --delimiter must not be empty")
			os.Exit(2)
		}
		var err error
		lines, err = readStdinRecords(os.Stdin, *delimiter)
		if err != nil {
			fmt.Fprintln(os.Stderr, "romanfmt: reading stdin:", err)
			os.Exit(1)
		}
	}

	if len(lines) == 0 {
		fmt.Fprintln(os.Stderr, "usage: romanfmt [--json] [--to-roman] [--strict] [--delimiter SEP] NUMERAL...")
		os.Exit(2)
	}

	exitCode := 0
	enc := json.NewEncoder(os.Stdout)
	for _, line := range lines {
		normalized, value, err := convert(line, *toRoman, *strict)
		r := result{Input: line}
		if err != nil {
			r.Error = err.Error()
			exitCode = 1
		} else {
			r.Normalized = normalized
			r.Value = value
		}

		if *jsonOutput {
			if err := enc.Encode(r); err != nil {
				fmt.Fprintln(os.Stderr, "romanfmt: encoding json:", err)
				os.Exit(1)
			}
			continue
		}

		if r.Error != "" {
			fmt.Printf("%s -> error: %s\n", r.Input, r.Error)
		} else {
			fmt.Printf("%s -> %s (%d)\n", r.Input, r.Normalized, r.Value)
		}
	}

	os.Exit(exitCode)
}

// convert dispatches to the Roman-to-decimal or decimal-to-Roman path
// depending on the --to-roman flag.
func convert(line string, toRoman, strict bool) (string, int, error) {
	if !toRoman {
		if strict {
			return NormalizeStrict(line)
		}
		return Normalize(line)
	}

	n, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return "", 0, fmt.Errorf("%q is not a decimal number", line)
	}
	roman, err := ToRoman(n)
	if err != nil {
		return "", 0, err
	}
	return roman, n, nil
}

// escapeReplacer expands the backslash escapes users are likely to type on a
// command line (e.g. --delimiter='\t') into their literal byte, since shells
// hand flag values to us unescaped otherwise.
var escapeReplacer = strings.NewReplacer(`\n`, "\n", `\t`, "\t", `\r`, "\r", `\0`, "\x00")

func readStdinRecords(r io.Reader, delimiter string) ([]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	delimiter = escapeReplacer.Replace(delimiter)
	var records []string
	for _, rec := range strings.Split(string(data), delimiter) {
		if rec == "" {
			continue
		}
		records = append(records, rec)
	}
	return records, nil
}
