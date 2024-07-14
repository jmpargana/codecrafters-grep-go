package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type re struct {
	predicates      map[string]func(rune) bool
	groupPredicates map[string]func([]byte) bool
	base            []byte
}

func (r *re) match() bool {
	m := map[rune]struct{}{}
	for _, b := range r.base {
		m[rune(b)] = struct{}{}
	}
	fmt.Println(r.predicates)
	fmt.Println(r.groupPredicates)
	for k := range m {
		for _, p := range r.predicates {
			if p(k) {
				return true
			}
		}
	}

	for _, p := range r.groupPredicates {
		if p(r.base) {
			return true
		}
	}

	return false
}

func (r *re) canMatch() bool {
	return len(r.predicates) > 0 || len(r.groupPredicates) > 0
}

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	pred := map[string]func(rune) bool{}
	groupPred := map[string]func([]byte) bool{}

	if pattern == "\\d" {
		pred["match digit"] = unicode.IsDigit
	}

	if pattern == "\\w" {
		pred["match digit"] = unicode.IsDigit
		pred["match alphanumeric"] = unicode.IsLetter
	}

	if strings.HasPrefix(pattern, "[") {
		positiveChars := strings.TrimSuffix(strings.TrimPrefix(pattern, "["), "]")
		if strings.HasPrefix(positiveChars, "^") {
			positiveChars = strings.TrimPrefix(positiveChars, "^")

			groupPred[fmt.Sprintf("match negative %s", positiveChars)] = func(s []byte) bool {
				sm := map[byte]struct{}{}

				for _, b := range s {
					sm[b] = struct{}{}
				}

				// forbidden map
				fm := map[byte]struct{}{}

				for _, b := range positiveChars {
					fm[byte(b)] = struct{}{}
				}

				for k := range sm {
					if _, ok := fm[k]; ok {
						return false
					}
				}

				return true
			}

		} else {
			for _, c := range positiveChars {
				pred[fmt.Sprintf("match positive rune %s", string(c))] = func(r rune) bool {
					return r == c
				}
			}
		}
	}

	r := re{base: line, predicates: pred, groupPredicates: groupPred}

	if r.canMatch() {
		return r.match(), nil
	}

	if utf8.RuneCountInString(pattern) != 1 && !r.canMatch() {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	ok := bytes.ContainsAny(line, pattern)
	return ok, nil
}
