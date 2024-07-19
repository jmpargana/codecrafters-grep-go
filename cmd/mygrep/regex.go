package main

type reTyp int
type carTyp int

const (
	char reTyp = iota
	digit
	alpha
	group
	begin
	end
	carMultiple
	carOptional
	carAny
	wildcard
)

const (
	optional carTyp = iota
	single
	multiple
	any
)

type RE struct {
	kind        reTyp
	value       rune
	subGroup    []RE
	negative    bool
	cardinality carTyp
}

func newSpec(kind reTyp) RE {
	return RE{kind, '*', nil, false, single}
}

func newChar(c byte) RE {
	return RE{char, rune(c), nil, false, single}
}

func newGroup(r []RE, neg bool) RE {
	return RE{group, '*', r, neg, single}
}

func parse(expr string) []RE {
	re := []RE{}

	for expr != "" {
		if expr[0] == '\\' {
			switch expr[1] {
			case 'd':
				re = append(re, newSpec(digit))
				expr = expr[2:]
			case 'w':
				re = append(re, newSpec(alpha))
				expr = expr[2:]
			default:
				re = append(re, newChar('\\'))
				expr = expr[2:]
			}
		} else if expr[0] == '[' {
			expr = expr[1:]
			negative := false
			if expr[0] == '^' {
				negative = true
				expr = expr[1:]
			}
			i := 1
			for ; i < len(expr); i++ {
				if expr[i] == ']' {
					break
				}
			}
			charGroup := parse(expr[:i])

			expr = expr[i+1:]
			re = append(re, newGroup(charGroup, negative))
		} else if expr[0] == '^' {
			re = append(re, newSpec(begin))
			expr = expr[1:]
		} else if expr[0] == '$' {
			re = append(re, newSpec(end))
			expr = expr[1:]
		} else if expr[0] == '+' {
			re = append(re, newSpec(carMultiple))
			expr = expr[1:]
		} else if expr[0] == '*' {
			re = append(re, newSpec(carAny))
			expr = expr[1:]
		} else if expr[0] == '?' {
			re = append(re, newSpec(carOptional))
			expr = expr[1:]
		} else if expr[0] == '.' {
			re = append(re, newSpec(wildcard))
			expr = expr[1:]
		} else {
			re = append(re, newChar(expr[0]))
			expr = expr[1:]
		}
	}

	return flattenCardinality(re)
}

func flattenCardinality(re []RE) []RE {
	for i := 0; i < len(re); i++ {
		switch re[i].kind {
		case carMultiple, carOptional, carAny:
			re[i-1].cardinality = getCardinality(re[i].kind)
			re = append(re[:i], re[i+1:]...)
		}
	}
	return re
}

func getCardinality(kind reTyp) carTyp {
	switch kind {
	case carMultiple:
		return multiple
	case carOptional:
		return optional
	case carAny:
		return any
	}
	return single
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func isNumeric(r byte) bool {
	return r >= '0' && r <= '9'
}

func isAlphaNumeric(r byte) bool {
	return isNumeric(r) || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func matchSingle(b byte, re RE) bool {
	switch re.kind {
	case char:
		return rune(b) == re.value
	case digit:
		return isNumeric(b)
	case alpha:
		return isAlphaNumeric(b)
	case wildcard:
		return true
	default:
		return false
	}
}

func match(expr, text string) bool {
	re := parse(expr)
	for i := 0; i < len(text); i++ {
		if matchRecursive(re, text, i) {
			return true
		}
	}
	return false
}

func matchRecursive(r []RE, text string, i int) bool {
	if len(r) == 0 {
		return i <= len(text)
	}

	re := r[0]

	switch re.kind {
	case begin:
		if i == 0 {
			return matchRecursive(r[1:], text, i)
		}
	case end:
		if i == len(text) {
			return true
		}
	case char, digit, alpha, wildcard:
		if i < len(text) && (matchSingle(text[i], re) || re.cardinality == optional) {
			if re.cardinality == optional {
				// FIXME: could need additional validation
				return matchRecursive(r[1:], text, i+1) || matchRecursive(r[1:], text, i)
			}
			if re.cardinality == multiple {
				return matchRecursive(r, text, i+1) || matchRecursive(r[1:], text, i+1)
			}
			return matchRecursive(r[1:], text, i+1)
		}
	case group:
		if re.negative {
			for _, opt := range re.subGroup {
				if matchSingle(text[i], opt) {
					return false
				}
			}
			return true
		} else {
			for _, opt := range re.subGroup {
				if matchSingle(text[i], opt) {
					if re.cardinality == multiple {
						return matchRecursive(r, text, i+1) || matchRecursive(r[1:], text, i+1)
					}
					return matchRecursive(r[1:], text, i+1)
				}
			}
		}
	}

	return false
}
