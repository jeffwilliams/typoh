package main

import (
	"bytes"
	"io"
	"testing"
)

func tokensEqual(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

func TestToker(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output []Token
	}{
		{
			name:  "simple",
			input: "<h>text</h>",
			output: []Token{
				Token{0, TokMetadata, "<h>"},
				Token{1, TokContent, "text"},
				Token{2, TokMetadata, "</h>"},
			},
		},
		{
			name:  "text only",
			input: "text",
			output: []Token{
				Token{0, TokContent, "text"},
			},
		},
		{
			name:  "tag only",
			input: "<tag>",
			output: []Token{
				Token{0, TokMetadata, "<tag>"},
			},
		},
		{
			name:  "open close",
			input: "<tag></tag>",
			output: []Token{
				Token{0, TokMetadata, "<tag>"},
				Token{1, TokMetadata, "</tag>"},
			},
		},
		{
			name:  "unclosed",
			input: "<tag>text is here",
			output: []Token{
				Token{0, TokMetadata, "<tag>"},
				Token{1, TokContent, "text is here"},
			},
		},
		{
			name:  "incomplete",
			input: "<tag",
			output: []Token{
				Token{0, TokMetadata, "<tag"},
			},
		},
		{
			name:  "a few tags",
			input: "<html><head>thing</head><body>This is some <i>text</i>.</html>",
			output: []Token{
				Token{0, TokMetadata, "<html>"},
				Token{1, TokMetadata, "<head>"},
				Token{2, TokContent, "thing"},
				Token{3, TokMetadata, "</head>"},
				Token{4, TokMetadata, "<body>"},
				Token{5, TokContent, "This is some "},
				Token{6, TokMetadata, "<i>"},
				Token{7, TokContent, "text"},
				Token{8, TokMetadata, "</i>"},
				Token{9, TokContent, "."},
				Token{10, TokMetadata, "</html>"},
			},
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			buf := bytes.NewBufferString(tc.input)
			toker := NewHtmlToker(buf)

			toks := []Token{}

			for {
				tok, err := toker.Next()
				if err != nil {
					if err == io.EOF {
						if tok != nil {
							toks = append(toks, *tok)
						}
					} else {
						t.Fatalf("Got error tokenizing: %v", err)
					}
					break
				}
				toks = append(toks, *tok)
			}

			if !tokensEqual(toks, tc.output) {
				t.Fatalf("Tokenizing was wrong for '%s'. Expected %v but got %v", tc.input, tc.output, toks)

			}

		})
	}
}
