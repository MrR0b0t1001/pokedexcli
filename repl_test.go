package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	// ...
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " pikachu bulbasaur chipotle  ",
			expected: []string{"pikachu", "bulbasaur", "chipotle"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Lengths don't match")
		}

		for idx := range actual {
			got := actual[idx]

			if !reflect.DeepEqual(got, c.expected[idx]) {
				t.Errorf("CleanInput(%q) = %v; want %v", c.input, got, c.expected)
			}

		}
	}
}
