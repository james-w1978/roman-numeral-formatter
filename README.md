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
