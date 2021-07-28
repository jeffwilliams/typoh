package main

import (
	"bytes"
	"fmt"
)

/*
typoh: Typography for HTML. Or a typo.
*/
/*
Replacements:
	-- = endash
	--- = emdash
	` = left curly
	' = right curly
	`` = left curlies
	'' = right curlies
	\ = escape the next character/token
	... = ellipses

Also figure out where a non-breaking space should go? Single underscore, with no space on each side only a non-_ letter?
*/

/*
Can use Butterick's idea of adding HTML soft hyphens to text using Liang's hyphenation algorithm from TeX
as implemented here: https://github.com/speedata/hyphenation
*/

type Toker interface {
	Next() (token *Token, done bool, err error)
}

type Token struct {
	index uint
	typ   TokenType
	text  string
}
type TokenType int

const (
	// TokMetadata are things that are not body text and so should not be formatted: HTML tags.
	TokMetadata = iota
	// TokContent is actual content
	TokContent
)

type Replacer struct {
	seq      []string
	repl     []string
	seqRunes [][]rune
	ndx      []int

	initialized bool
	input       []rune
	result      bytes.Buffer
	holdback    []rune
}

func (r *Replacer) Add(sequence, replacement string) {
	r.init()

	r.seqRunes = append(r.seqRunes, []rune(sequence))
	r.seq = append(r.seq, sequence)
	r.repl = append(r.repl, replacement)
	r.ndx = append(r.ndx, 0)

}

func (r *Replacer) init() {
	if r.initialized {
		return
	}
	r.seq = make([]string, 0, 10)
	r.repl = make([]string, 0, 10)
	r.seqRunes = make([][]rune, 0, 10)
	r.ndx = make([]int, 0, 10)
	r.holdback = make([]rune, 0, 10)

	r.initialized = true
}

func (r *Replacer) Replace(text string) string {
	r.init()

	r.input = []rune(text)
	r.result.Reset()

	for _, rn := range r.input {
		r.advance(rn)
	}

	// Flush out the holdback.
	r.result.WriteString(string(r.holdback))

	return r.result.String()
}

func (r *Replacer) advance(rn rune) {
	fmt.Printf("advancing on rune %c\n", rn)

	matchIndex := -1

	longestPotentialMatchLenBeforeRune := 0

	for i := range r.ndx {
		if r.ndx[i]+1 > longestPotentialMatchLenBeforeRune {
			longestPotentialMatchLenBeforeRune = r.ndx[i] + 1
		}
	}

	for i, seq := range r.seqRunes {
		fmt.Printf(" check sequence %s\n", string(seq))

		if seq[r.ndx[i]] == rn {
			r.ndx[i]++
			fmt.Printf(" sequence %s: moved ahead to index %d\n", string(seq), r.ndx[i])
		} else {
			r.ndx[i] = 0
			fmt.Printf(" sequence %s: reset ahead to index %d\n", string(seq), r.ndx[i])
		}

		if r.ndx[i] == len(seq) {
			matchIndex = i
			r.ndx[i] = 0
		}
	}

	if matchIndex >= 0 {
		fmt.Printf(" sequence %s matched. Replacing with %s\n", string(r.seq[matchIndex]), r.repl[matchIndex])

		// We have a match. Output the replacement instead of everything we have in the holdback.
		// Clear the holdback, and reset all possible matches
		r.result.WriteString(r.repl[matchIndex])
		r.holdback = r.holdback[:0]
		for i := range r.ndx {
			r.ndx[i] = 0
		}

		return
	}

	// No match
	r.holdback = append(r.holdback, rn)

	longestPotentialMatchLenAfterRune := 0

	for i := range r.seqRunes {
		if r.ndx[i]+1 > longestPotentialMatchLenAfterRune {
			longestPotentialMatchLenAfterRune = r.ndx[i] + 1
		}
	}

	if len(r.seq) > 0 && longestPotentialMatchLenAfterRune < longestPotentialMatchLenBeforeRune+1 {
		fmt.Printf(" longest potential match failed to match. longest before=%d, longest after=%d, holdback = %s\n",
			longestPotentialMatchLenBeforeRune, longestPotentialMatchLenAfterRune, string(r.holdback))

		// Our previous longest potential matches turned out not to match.
		// We can now write out
		// the runes from the holdback that we were holding back, but only up to the
		// point of the _next_ longest match that can still match.
		n := longestPotentialMatchLenBeforeRune + 1 - longestPotentialMatchLenAfterRune
		// write out n runes and move the runes after the first n from holdback to the front of holdback.
		r.result.WriteString(string(r.holdback[0:n]))
		copy(r.holdback, r.holdback[n:])
		r.holdback = r.holdback[0 : len(r.holdback)-n]
	}
}

func main() {

}
