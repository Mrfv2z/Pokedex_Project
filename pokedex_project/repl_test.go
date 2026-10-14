package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " SUP DUDE ",
			expected: []string{"sup", "dude"},
		},
		{
			input:    "  NO MERCY!  ",
			expected: []string{"no", "mercy!"},
		},
	}

	for _, tc := range cases {
		if got := cleanInput(tc.input); !reflect.DeepEqual(got, tc.expected) {
			t.Fatalf("cleanInput(%q) = %#v; want %#v", tc.input, got, tc.expected)
		}
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice
		// if they don't match, use t.Errorf and continue to the next case
		if len(actual) != len(c.expected) {
			t.Errorf("size mismatch")
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("words don't match")
				fmt.Println("FAIL!")
			}
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
		}
	}
}
