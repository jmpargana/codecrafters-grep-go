package main

import (
	"reflect"
	"testing"
)

func TestBase(t *testing.T) {
	tt := []struct {
		expr     string
		text     string
		expected bool
	}{
		{"abc", "abc", true},
		{"abc", "adc", false},
		{"abc", "abd", false},
		{"abc", "abb", false},
		{"abc", "bbc", false},
		{"a", "bbac", true},
		{"aaa", "aa", false},
		{"abc", "aaabccc", true},
		{`a\dc`, "a1c", true},
		{`a\dc`, "a9c", true},
		{`a\dc`, "a0c", true},
		{`a\wc`, "a0c", true},
		{`a\wc`, "abc", true},
		{`a\dc`, "abc", false},
		{`[abc]`, "a", true},
		{`[abc]`, "tttattt", true},
		{`[abc]`, "ttta", true},
		{`[abc]`, "attt", true},
		{`[abc]`, "tttt", false},
		{`[abc]`, "b", true},
		{`[abc]`, "c", true},
		{`[^abc]`, "c", false},
		{`[^ban]`, "banana", false},
		{`[^xyz]`, "banana", true},
		{`a[^xyz]b`, "abc", true},
		{`a[^xyz]b`, "axc", false},
		{`a[^xyz]b`, "axc", false},
		{`^log`, "log", true},
		{`^log`, "llog", false},
		{`dog$`, "dog", true},
		{`dog$`, "dogg", false},
		{`ca+ts`, "cats", true},
		{`ca+ts`, "caaaats", true},
		{`ca+ts`, "caabats", false},
		{`c[abc]+ts`, "caabats", true},
		{`\d+`, "1123", true},
		{`\d+`, "1", true},
		{`^\d+`, "a1", false},
		{`\d?a`, "a", true},
		{`\d?a`, "9a", true},
		{`^\d?a`, "90a", false},
		{`\d?a`, "90a", true},
		{`^[abc]?d`, "ad", true},
		{`^[abc]?d`, "bd", true},
		{`^[abc]?d`, "bbd", false},
	}

	for _, tc := range tt {
		received := match(tc.expr, tc.text)
		if tc.expected != received {
			t.Errorf("invalid result for: %s to be matched in %s\n", tc.expr, tc.text)
		}
	}
}

func TestParse(t *testing.T) {
	tt := []struct {
		expr     string
		expected []RE
	}{
		{"abc", []RE{
			newChar('a'),
			newChar('b'),
			newChar('c'),
		}},
		{`a\dc`, []RE{
			newChar('a'),
			newSpec(digit),
			newChar('c'),
		}},
		{`a\d`, []RE{
			newChar('a'),
			newSpec(digit),
		}},
		{`a\d\dc`, []RE{
			newChar('a'),
			newSpec(digit),
			newSpec(digit),
			newChar('c'),
		}},
		{`\da\d\dc`, []RE{
			newSpec(digit),
			newChar('a'),
			newSpec(digit),
			newSpec(digit),
			newChar('c'),
		}},
		{`a\wc`, []RE{
			newChar('a'),
			newSpec(alpha),
			newChar('c'),
		}},
		{`a[abc]c`, []RE{
			newChar('a'),
			newGroup([]RE{
				newChar('a'),
				newChar('b'),
				newChar('c'),
			}, false),
			newChar('c'),
		}},
		{`a[^abc]c`, []RE{
			newChar('a'),
			newGroup([]RE{
				newChar('a'),
				newChar('b'),
				newChar('c'),
			}, true),
			newChar('c'),
		}},
		{`a[ab\db]c`, []RE{
			newChar('a'),
			newGroup([]RE{
				newChar('a'),
				newChar('b'),
				newSpec(digit),
				newChar('b'),
			}, false),
			newChar('c'),
		}},
		{`a[ab\d]c`, []RE{
			newChar('a'),
			newGroup([]RE{
				newChar('a'),
				newChar('b'),
				newSpec(digit),
			}, false),
			newChar('c'),
		}},
		{`\d\\d\\d`, []RE{
			newSpec(digit),
			newChar('\\'),
			newChar('d'),
			newChar('\\'),
			newChar('d'),
		}},
		{`^log`, []RE{
			newSpec(begin),
			newChar('l'),
			newChar('o'),
			newChar('g'),
		}},
		{`dog$`, []RE{
			newChar('d'),
			newChar('o'),
			newChar('g'),
			newSpec(end),
		}},
		{`ca+ts`, []RE{
			newChar('c'),
			{char, 'a', nil, false, multiple},
			newChar('t'),
			newChar('s'),
		}},
		{`[^abc]+e`, []RE{
			{group, '*', []RE{
				newChar('a'),
				newChar('b'),
				newChar('c'),
			}, true, multiple},
			newChar('e'),
		}},
		{`\d+`, []RE{
			{digit, '*', nil, false, multiple},
		}},
		{`a?`, []RE{
			{char, 'a', nil, false, optional},
		}},
		{`\d?a`, []RE{
			{digit, '*', nil, false, optional},
			{char, 'a', nil, false, single},
		}},
	}
	for _, tc := range tt {
		received := parse(tc.expr)
		if !reflect.DeepEqual(received, tc.expected) {
			t.Errorf("\n\ncase: %s\nreceived:\t\t%v\nwanted:\t\t\t%v\n", tc.expr, received, tc.expected)
		}
	}
}
