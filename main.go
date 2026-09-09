package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
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
	flag.Parse()

	lines := flag.Args()
	if len(lines) == 0 {
		var err error
		lines, err = readStdinLines()
		if err != nil {
			fmt.Fprintln(os.Stderr, "romanfmt: reading stdin:", err)
			os.Exit(1)
		}
	}

	if len(lines) == 0 {
		fmt.Fprintln(os.Stderr, "usage: romanfmt [--json] [--to-roman] NUMERAL...")
		os.Exit(2)
	}

	exitCode := 0
	enc := json.NewEncoder(os.Stdout)
	for _, line := range lines {
		normalized, value, err := convert(line, *toRoman)
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
func convert(line string, toRoman bool) (string, int, error) {
	if !toRoman {
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

func readStdinLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}
