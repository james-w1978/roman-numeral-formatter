# roman-numeral-formatter

Roman numerals show up in the wild in every case and spelling variant you can
imagine: lowercase (`xiv`), stray whitespace (`M CM XC IV`), and non-canonical
repetition that's technically readable but not "correct" (`IIII` instead of
`IV`). This is a small CLI that takes any of that and normalizes it to the
one true canonical form, plus tells you the integer value.

It also parses leniently on the way in: anything with a smaller value before
a larger one is treated as subtractive, everything else is additive. That's
enough to make sense of sloppy input without a strict grammar.

## Usage

Build it:

```
go build -o romanfmt .
```

Pass numerals as arguments:

```
$ ./romanfmt xiv "  MCMxciv  " IIII
xiv -> XIV (14)
  MCMxciv   -> MCMXCIV (1994)
IIII -> IV (4)
```

Or pipe them in, one per line:

```
$ printf "xl\nix\n" | ./romanfmt
xl -> XL (40)
ix -> IX (9)
```

Bad input reports an error instead of a bogus result, and the process exits
non-zero if any input failed:

```
$ ./romanfmt XIIB
XIIB -> error: invalid character 'B' in "XIIB"
```

## Converting the other direction

Pass `--to-roman` to go from decimal to Roman numerals instead:

```
$ ./romanfmt --to-roman 1994 4
1994 -> MCMXCIV (1994)
4 -> IV (4)
```

It still rejects out-of-range values and non-numeric input:

```
$ ./romanfmt --to-roman 4000 abc
4000 -> error: value 4000 out of range 1-3999
abc -> error: "abc" is not a decimal number
```

## Custom stdin delimiters

By default, piped input is split on newlines. Pass `--delimiter` to split on
something else, which is handy for comma-separated or otherwise packed input:

```
$ printf "xiv,IIII,xl" | ./romanfmt --delimiter ","
xiv -> XIV (14)
IIII -> IV (4)
xl -> XL (40)
```

Common escapes (`\n`, `\t`, `\r`, `\0`) are recognized so you can pass them as
plain text on the command line instead of a literal control character.

## Strict mode

By default sloppy repetition like `IIII` is corrected rather than rejected.
Pass `--strict` to reject anything that isn't already in canonical form:

```
$ ./romanfmt --strict IIII XIV
IIII -> error: "IIII" is not canonical (expected "IV")
XIV -> XIV (14)
```

## JSON output

Pass `--json` to get one JSON object per line instead, which is easier to
pipe into other tools:

```
$ ./romanfmt --json xiv XIIB
{"input":"xiv","normalized":"XIV","value":14}
{"input":"XIIB","error":"invalid character 'B' in \"XIIB\""}
```

## Range

Standard Roman numerals only cover 1 to 3999 without inventing extra
notation (vinculum, apostrophus, etc.), so values outside that range are
rejected rather than guessed at.

## License

MIT, see LICENSE.
