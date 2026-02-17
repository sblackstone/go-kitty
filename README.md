# go-kitty

A tiny cat toy for your terminal (ncurses-ish via tcell). It’s basically a little screen full of slithery distractions meant for curious paws.

## What it is
- A terminal animation playground
- A cat toy for ncurses-style terminals
- A reason to leave your laptop open for “cat QA”

## Run
- `go run .`

## CLI options
Tune how many critters appear and their initial spawn delays.

Examples:
- `go run . --snakes 3 --snake-max-len 14`
- `go run . --strings 1 --string-min-len 12 --string-max-len 28`
- `go run . --butterflies 2 --butterfly-initial-delay-max 120`
- `go run . --lasers 2 --laser-initial-delay-max 60`

Flags:
- `--snakes` (default: 2)
- `--snake-max-len` (default: 10)
- `--snake-initial-delay-max` (default: 40)
- `--strings` (default: 2)
- `--string-min-len` (default: 18)
- `--string-max-len` (default: 36)
- `--string-initial-delay-max` (default: 40)
- `--butterflies` (default: 1)
- `--butterfly-initial-delay-max` (default: 80)
- `--lasers` (default: 1)
- `--laser-initial-delay-max` (default: 80)

## Disclaimer
Not responsible for unexpected pounces, keyboard naps, or the sudden disappearance of your cursor.
