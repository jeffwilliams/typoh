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
	=> = bullet
	1/2 = 1/2 symbol
	1/4 = 1/4 symbol
	3/4 = 3/4 symbol
	^N = superscript N, where N is a digit between 0-9 inclusive
	_N = subscript N, where N is a digit between 0-9 inclusive
	N^o = №

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
	r.Add("=>", "•")
	r.Add("N^o", "№")
	r.Add("1/2", "½")
	r.Add("1/4", "¼")
	r.Add("3/4", "¾")
	r.Add("^0", "⁰")
	r.Add("^1", "¹")
	r.Add("^2", "²")
	r.Add("^3", "³")
	r.Add("^4", "⁴")
	r.Add("^5", "⁵")
	r.Add("^6", "⁶")
	r.Add("^7", "⁷")
	r.Add("^8", "⁸")
	r.Add("^9", "⁹")
	r.Add("_0", "₀")
	r.Add("_1", "₁")
	r.Add("_2", "₂")
	r.Add("_3", "₃")
	r.Add("_4", "₄")
	r.Add("_5", "₅")
	r.Add("_6", "₆")
	r.Add("_7", "₇")
	r.Add("_8", "₈")
	r.Add("_9", "₉")
	
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
