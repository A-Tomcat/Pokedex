package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " bring me the   horizon",
			expected: []string{"bring", "me", "the", "horizon"},
		},
		{
			input:    "stupidity never won a price",
			expected: []string{"stupidity", "never", "won", "a", "price"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length Error: actual : %v, expected %v", len(actual), len(c.expected))
			t.Fail()
		}
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			//Check each word in the slice
			//if they don't match, use t.Errorf to print an error message
			//and fail the test
			if word != expectedWord {
				t.Errorf("expected Word: %s, does not match actual word: %s", expectedWord, word)
				t.Fail()
			}
		}
	}

}
