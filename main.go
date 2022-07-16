package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
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
	_ = nonbreaking space
	~ = optional hyphen? or _ See unicode: https://www.compart.com/en/unicode/U+00A0

Also figure out where a non-breaking space should go? Single underscore, with no space on each side only a non-_ letter?
*/

/*
Can use Butterick's idea of adding HTML soft hyphens to text using Liang's hyphenation algorithm from TeX
as implemented here: https://github.com/speedata/hyphenation
*/

func main() {
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

	r.Add("---", "—")
	r.Add(`\-`, "-")
	r.Add("``", "“")
	r.Add("''", "”")
	r.Add("...", "…")
	r.Add(`\.`, ".")
	r.Add("_", " ")
	r.Add("~", "­")

	return r
}

func (t *Typographer) secondPassReplacer() *Replacer {
	r := &Replacer{}

	r.Add(`\'`, "'")
	r.Add("--", "–")
	r.Add("`", "‘")
	r.Add("\\`", "`")
	r.Add("'", "’")

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
