package main

import (
	"testing"
)

func TestReplacer(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		output       string
		replacements [][]string
	}{

		{
			name:   "no replacements",
			input:  "abcd",
			output: "abcd",
		},
		{
			name:   "empty input and no replacements",
			input:  "",
			output: "",
		},
		{
			name:  "a to x",
			input: "this is a test",
			replacements: [][]string{
				[]string{"a", "x"},
			},
			output: "this is x test",
		},
		{
			name:  "a to x, t to y",
			input: "this is a test",
			replacements: [][]string{
				[]string{"a", "x"},
				[]string{"t", "y"},
			},
			output: "yhis is x yesy",
		},
		{
			name:  "abc, bde",
			input: "abde",
			replacements: [][]string{
				[]string{"abc", "x"},
				[]string{"bde", "y"},
			},
			output: "ay",
		},
		{
			name:  "abc, bde2",
			input: "matchabde",
			replacements: [][]string{
				[]string{"abc", "x"},
				[]string{"bde", "y"},
			},
			output: "matchay",
		},
		{
			name:  "abc, bde3",
			input: "matchabdes",
			replacements: [][]string{
				[]string{"abc", "x"},
				[]string{"bde", "y"},
			},
			output: "matchays",
		},
		{
			name:  "abc, bde",
			input: "match abdeabc",
			replacements: [][]string{
				[]string{"abc", "x"},
				[]string{"bde", "y"},
			},
			output: "match ayx",
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			var replacer Replacer

			for _, r := range tc.replacements {
				replacer.Add(r[0], r[1])
			}

			output := replacer.Replace(tc.input)
			if tc.output != output {

				t.Fatalf("Replacement was wrong for '%s'. Expected %v but got %v", tc.input, tc.output, output)

			}

		})
	}
}
