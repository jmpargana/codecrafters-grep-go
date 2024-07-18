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
			{char, 'a', nil, false},
			{char, 'b', nil, false},
			{char, 'c', nil, false},
		}},
		{`a\dc`, []RE{
			{char, 'a', nil, false},
			{digit, '*', nil, false},
			{char, 'c', nil, false},
		}},
		{`a\d`, []RE{
			{char, 'a', nil, false},
			{digit, '*', nil, false},
		}},
		{`a\d\dc`, []RE{
			{char, 'a', nil, false},
			{digit, '*', nil, false},
			{digit, '*', nil, false},
			{char, 'c', nil, false},
		}},
		{`\da\d\dc`, []RE{
			{digit, '*', nil, false},
			{char, 'a', nil, false},
			{digit, '*', nil, false},
			{digit, '*', nil, false},
			{char, 'c', nil, false},
		}},
		{`a\wc`, []RE{
			{char, 'a', nil, false},
			{alpha, '*', nil, false},
			{char, 'c', nil, false},
		}},
		{`a[abc]c`, []RE{
			{char, 'a', nil, false},
			{group, '*', []RE{
				{char, 'a', nil, false},
				{char, 'b', nil, false},
				{char, 'c', nil, false},
			}, false},
			{char, 'c', nil, false},
		}},
		{`a[^abc]c`, []RE{
			{char, 'a', nil, false},
			{group, '*', []RE{
				{char, 'a', nil, false},
				{char, 'b', nil, false},
				{char, 'c', nil, false},
			}, true},
			{char, 'c', nil, false},
		}},
		{`a[ab\db]c`, []RE{
			{char, 'a', nil, false},
			{group, '*', []RE{
				{char, 'a', nil, false},
				{char, 'b', nil, false},
				{digit, '*', nil, false},
				{char, 'b', nil, false},
			}, false},
			{char, 'c', nil, false},
		}},
		{`a[ab\d]c`, []RE{
			{char, 'a', nil, false},
			{group, '*', []RE{
				{char, 'a', nil, false},
				{char, 'b', nil, false},
				{digit, '*', nil, false},
			}, false},
			{char, 'c', nil, false},
		}},
		{`\d\\d\\d`, []RE{
			{digit, '*', nil, false},
			{char, '\\', nil, false},
			{char, 'd', nil, false},
			{char, '\\', nil, false},
			{char, 'd', nil, false},
		}},
		{`^log`, []RE{
			{begin, '*', nil, false},
			{char, 'l', nil, false},
			{char, 'o', nil, false},
			{char, 'g', nil, false},
		}},
		{`dog$`, []RE{
			{char, 'd', nil, false},
			{char, 'o', nil, false},
			{char, 'g', nil, false},
			{end, '*', nil, false},
		}},
	}
	for _, tc := range tt {
		received := parse(tc.expr)
		if !reflect.DeepEqual(received, tc.expected) {
			t.Errorf("case: %s, received: %v, wanted %v\n", tc.expr, received, tc.expected)
		}
	}
}
