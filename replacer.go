package main

import (
	"bytes"
	"fmt"
)

type Replacer struct {
	seq      []string
	repl     []string
	seqRunes [][]rune
	ndx      []int
	debug    bool

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

// Replace matches in text. If two replacements would both match at the same time, the first one Added is used.
// TODO: test this assertion
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
	r.dbg("advancing on rune %c\n", rn)

	longestPotentialMatchLenBeforeRune := r.longestPotentialMatchLen()

	matchIndex := r.advanceEachSequence(rn)

	if matchIndex >= 0 {
		r.dbg(" sequence %s matched. Replacing with %s\n", string(r.seq[matchIndex]), r.repl[matchIndex])

		// We have a match. Output the replacement instead of the last n characters we have in the holdback.
		r.writeFirstNRunesOfHoldback(len(r.holdback) - len(r.seqRunes[matchIndex]) + 1)
		r.result.WriteString(r.repl[matchIndex])
		r.clearHoldback()
		//	r.holdback = r.holdback[0 : len(r.seqRunes[matchIndex])-1]
		r.resetIndexes()

		return
	}

	// No match
	r.appendToHoldback(rn)

	longestPotentialMatchLenAfterRune := r.longestPotentialMatchLen()

	if len(r.seq) > 0 && longestPotentialMatchLenAfterRune < longestPotentialMatchLenBeforeRune+1 {
		r.dbg(" longest potential match failed to match. longest before=%d, longest after=%d, holdback = %s\n",
			longestPotentialMatchLenBeforeRune, longestPotentialMatchLenAfterRune, string(r.holdback))

		// Our previous longest potential matches turned out not to match.
		// We can now write out
		// the runes from the holdback that we were holding back, but only up to the
		// point of the _next_ longest match that can still match.
		n := longestPotentialMatchLenBeforeRune + 1 - longestPotentialMatchLenAfterRune
		// write out n runes and move the runes after the first n from holdback to the front of holdback.
		r.writeFirstNRunesOfHoldback(n)
	}
}

func (r *Replacer) resetIndexes() {
	for i := range r.ndx {
		r.ndx[i] = 0
	}
}

func (r *Replacer) clearHoldback() {
	r.holdback = r.holdback[:0]
}

func (r *Replacer) appendToHoldback(rn rune) {
	r.holdback = append(r.holdback, rn)
}

func (r *Replacer) writeFirstNRunesOfHoldback(n int) {
	r.result.WriteString(string(r.holdback[0:n]))
	copy(r.holdback, r.holdback[n:])
	r.holdback = r.holdback[0 : len(r.holdback)-n]
}

func (r Replacer) advanceEachSequence(rn rune) (matchIndex int) {
	matchIndex = -1

	for i, seq := range r.seqRunes {
		r.dbg(" check sequence %s\n", string(seq))

		if seq[r.ndx[i]] == rn {
			r.ndx[i]++
			r.dbg(" sequence %s: moved ahead to index %d\n", string(seq), r.ndx[i])
		} else {
			r.ndx[i] = 0
			r.dbg(" sequence %s: reset ahead to index %d\n", string(seq), r.ndx[i])
		}

		if r.ndx[i] == len(seq) && matchIndex == -1 {
			matchIndex = i
			r.ndx[i] = 0
		}
	}
	return
}

func (r Replacer) longestPotentialMatchLen() int {
	n := 0
	for i := range r.seqRunes {
		if r.ndx[i]+1 > n {
			n = r.ndx[i] + 1
		}
	}
	return n
}

func (r Replacer) dbg(format string, args ...interface{}) {
	if !r.debug {
		return
	}

	fmt.Printf(format, args...)
}
