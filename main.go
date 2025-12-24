/*
typoh: Typography for HTML. Or a typo.

Replacements:

	-- = endash
	--- = emdash
	` = left curly
	' = right curly
	`` = left curlies
	'' = right curlies
	\ = escape the next character/token
	... = ellipses
	_ = nonbreaking space
	~ = optional hyphen? or _ See unicode: https://www.compart.com/en/unicode/U+00A0
	=> = bullet
	1/2 = 1/2 symbol
	1/4 = 1/4 symbol
	3/4 = 3/4 symbol
	1/3 = 1/3 symbol
	2/3 = 2/3 symbol
	^N = superscript N, where N is a digit between 0-9 inclusive
	_N = subscript N, where N is a digit between 0-9 inclusive
	N^o = №
	:check: = "✓"
	:rightarrow: = →  (and also left, up, down)
	:mult: = "×"

*/
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)


/*
Can use Butterick's idea of adding HTML soft hyphens to text using Liang's hyphenation algorithm from TeX
as implemented here: https://github.com/speedata/hyphenation
*/

var (
	firstPassReplacements = []Replacement{
		{"---", "—"},
		{`\-`, "-"},
		{"``", "“"},
		{"''", "”"},
		{"...", "…"},
		{`\.`, "."},
		{"~", "­"},
		{"=>", "•"},
		{"N^o", "№"},
		{"1/2", "½"},
		{"1/4", "¼"},
		{"3/4", "¾"},
		{"1/3", "⅓"},
		{"2/3", "⅔"},
		{"^0", "⁰"},
		{"^1", "¹"},
		{"^2", "²"},
		{"^3", "³"},
		{"^4", "⁴"},
		{"^5", "⁵"},
		{"^6", "⁶"},
		{"^7", "⁷"},
		{"^8", "⁸"},
		{"^9", "⁹"},
		{"_0", "₀"},
		{"_1", "₁"},
		{"_2", "₂"},
		{"_3", "₃"},
		{"_4", "₄"},
		{"_5", "₅"},
		{"_6", "₆"},
		{"_7", "₇"},
		{"_8", "₈"},
		{"_9", "₉"},
		{":check:", "✓"},
		{":rightarrow:", "→"},
		{":leftarrow:", "←"},
		{":uparrow:", "↑"},
		{":downarrow:", "↓"},
		{":mult:", "×"},
	}

	secondPassReplacements = []Replacement{
		{"_", " "},
		{`\'`, "'"},
		{"--", "–"},
		{"`", "‘"},
		{"\\`", "`"},
		{"'", "’"},
	}
)

func main() {
	pflag.Usage = usage
	pflag.Parse()

	input, err := inputStream()
	if err != nil {
		return
	}

	var typo Typographer
	typo.ReplaceMarkers(input, os.Stdout)
}

func inputStream() (io.Reader, error) {
	if pflag.NArg() < 1 {
		return os.Stdin, nil
	}

	fname := pflag.Arg(0)
	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("Error opening %s: %v\n", fname, err)
	}
	return file, nil
}

type Typographer struct {
	replacers []*Replacer
	output    io.Writer
	toker     Toker
}

func (t *Typographer) ReplaceMarkers(input io.Reader, output io.Writer) error {
	t.output = output
	t.initReplacers()
	t.initTokenizer(input)
	return t.replaceMarkers()
}

func (t *Typographer) initReplacers() {
	if t.replacers == nil {
		t.replacers = []*Replacer{
			t.firstPassReplacer(),
			t.secondPassReplacer(),
		}
	}
}

func (t *Typographer) firstPassReplacer() *Replacer {
	r := &Replacer{}

	for _, repl := range firstPassReplacements {
		r.Add(repl.from, repl.to)
	}

	return r
}

func (t *Typographer) secondPassReplacer() *Replacer {
	r := &Replacer{}

	for _, repl := range secondPassReplacements {
		r.Add(repl.from, repl.to)
	}

	return r
}

func (t *Typographer) initTokenizer(input io.Reader) {
	t.toker = NewHtmlToker(input)
}

func (t *Typographer) replaceMarkers() error {
	for {
		tok, err := t.toker.Next()

		if tok != nil {
			t.replaceMarkersIn(tok)
			fmt.Fprintf(t.output, tok.text)
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("Got error reading and tokenizing input: %v", err)
		}
	}
	return nil
}

func (t *Typographer) replaceMarkersIn(tok *Token) {
	if tok.typ == TokContent {
		for _, r := range t.replacers {
			tok.text = r.Replace(tok.text)
		}
	}
}

type Replacement struct {
	from, to string
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s < <file>\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Substitutions: \n")
	for _, repl := range firstPassReplacements {
		fmt.Fprintf(os.Stderr, "  %s --> %s\n", repl.from, repl.to)
	}
	for _, repl := range secondPassReplacements {
		fmt.Fprintf(os.Stderr, "  %s --> %s\n", repl.from, repl.to)
	}
	//pflag.PrintDefaults()
}
