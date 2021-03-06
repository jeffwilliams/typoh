package main

import (
	"bufio"
	"bytes"
	"io"
)

type Toker interface {
	// Next returns io.EOF when done. It may also return a token at the same time.
	Next() (token *Token, err error)
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

type HtmlToker struct {
	rdr   *bufio.Reader
	inTag bool
	index uint
	text  bytes.Buffer
}

func NewHtmlToker(r io.Reader) *HtmlToker {
	return &HtmlToker{
		rdr: bufio.NewReader(r),
	}
}

func (h *HtmlToker) Next() (token *Token, err error) {

	for {
		var r rune
		r, _, err = h.rdr.ReadRune()

		if err != nil {
			if !h.currentTokenIsEmpty() {
				token = h.mkCurrentToken()
			}
			return
		}

		if h.inTag {
			h.text.WriteRune(r)
			if r == '>' {
				token = h.mkCurrentToken()
				h.inTag = false
				return
			}
		} else {
			if r == '<' {
				empty := h.currentTokenIsEmpty()
				if !empty {
					token = h.mkCurrentToken()
				}
				h.inTag = true

				if !empty {
					h.text.WriteRune(r)
					return
				}
			}
			h.text.WriteRune(r)
		}
	}
	return
}

func (h *HtmlToker) currentTokenIsEmpty() bool {
	return h.text.Len() == 0
}

func (h *HtmlToker) mkCurrentToken() (tok *Token) {
	if h.inTag {
		tok = h.mkToken(TokMetadata)
	} else {
		tok = h.mkToken(TokContent)
	}
	h.index++
	return
}

func (h *HtmlToker) mkToken(typ TokenType) *Token {
	txt := h.text.String()
	h.text.Reset()
	return &Token{
		index: h.index,
		typ:   typ,
		text:  txt,
	}
}
