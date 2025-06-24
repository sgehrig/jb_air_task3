package cli

import (
	"reflect"
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		input string
		want  []string
		fail  bool
	}{
		// Simple
		{"foo bar baz", []string{"foo", "bar", "baz"}, false},
		// Quoted
		{"foo 'bar baz' qux", []string{"foo", "bar baz", "qux"}, false},
		{"foo \"bar baz\" qux", []string{"foo", "bar baz", "qux"}, false},
		// Escaped quotes inside
		{`foo 'bar\'s' qux`, []string{"foo", "bar's", "qux"}, false},
		{`foo "bar\"s" qux`, []string{"foo", "bar\"s", "qux"}, false},
		// Nested quotes (escaped)
		{"foo 'bar \"baz\"' qux", []string{"foo", "bar \"baz\"", "qux"}, false},
		{"foo \"bar 'baz'\" qux", []string{"foo", "bar 'baz'", "qux"}, false},
		// Mixed whitespace
		{"foo\tbar  'baz qux'", []string{"foo", "bar", "baz qux"}, false},
		// Leading/trailing whitespace
		{"  foo  bar  ", []string{"foo", "bar"}, false},
		// Only whitespace
		{"   \t  ", nil, false},
		// Empty
		{"", nil, false},
		// Unclosed quote
		{"foo 'bar", nil, true},
		{"foo \"bar", nil, true},
		// Unfinished escape
		{"foo bar\\", nil, true},
		// Escaped space
		{"foo bar\\ baz", []string{"foo", "bar baz"}, false},
		// Escaped quote at end (should fail: unclosed quote)
		{`foo 'bar\'`, nil, true},
		// Quotes inside quotes
		{"Test 'ass'", []string{"Test", "ass"}, false},
		{"Test \"ABC\"", []string{"Test", "ABC"}, false},
		{"'Test \"ABC\"'", []string{"Test \"ABC\""}, false},
		{"\"Test 'ass'\"", []string{"Test 'ass'"}, false},
	}
	for _, tt := range tests {
		got, err := ParseCommandLine(tt.input)
		if tt.fail {
			if err == nil {
				t.Errorf("expected error for input %q, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error for input %q: %v", tt.input, err)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParseCommandLine(%q) = %#v, want %#v", tt.input, got, tt.want)
		}
	}
}
