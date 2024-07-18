package main

type reTyp int

const (
	char reTyp = iota
	digit
	alpha
	group
)

type RE struct {
	kind     reTyp
	value    rune
	subGroup []RE
	negative bool
}

func parse(expr string) []RE {
	re := []RE{}

	for expr != "" {
		if expr[0] == '\\' {
			if len(expr) < 2 {
				break
			}
			switch expr[1] {
			case 'd':
				re = append(re, RE{digit, '*', nil, false})
				expr = expr[2:]
			case 'w':
				re = append(re, RE{alpha, '*', nil, false})
				expr = expr[2:]
			default:
				re = append(re, RE{char, '\\', nil, false})
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
			re = append(re, RE{group, '*', charGroup, negative})
		} else {
			re = append(re, RE{char, rune(expr[0]), nil, false})
			expr = expr[1:]
		}
	}

	return re
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
	case char, digit, alpha:
		if i < len(text) && matchSingle(text[i], re) {
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
					return matchRecursive(r[1:], text, i+1)
				}
			}
		}
	}

	return false
}
