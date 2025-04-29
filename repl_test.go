package main

import "testing"

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
			input: " ",
			expected: []string{},
		},
		{
			input: "aloHa woRLd",
			expected: []string{"aloha","world"},
		},
		{
			input: "     asDf!#$ $ $$T   ",
			expected: []string{"asdf!#$","$","$$t"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("Slice lengths don't match for: %v", c.input)
			t.Fail()
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Output word didn't match: ACTUAL = %v, EXPECTED = %v", word, expectedWord)
				t.Fail()
			}
		}
	}
}