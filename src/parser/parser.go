package parser

import (
	"bufio"
	"bytes"
	"io"
)

type AttributeTemplate struct {
	Name     Attr
	Template []Attr
}

func (a *AttributeTemplate) Parse(scanrrr *bufio.Scanner, w *bytes.Buffer) {
	boiler := &AttributeTemplate{
		Name: "BOILER",
		Template: []Attr{
			"HTML", "BODY",
			"DIV", "@PARAM", "TITLE", "END",
			"DIV", "@PARAM", "INSIDE", "END",
			"END", "END"},
	}
}

func ParseCum(in io.Reader) (out io.Reader) {
	w := new(bytes.Buffer)
	out = bufio.NewReader(w)

	scanrrr := bufio.NewScanner(in)

	for scanrrr.Scan() {
		// word := Token(scanrrr.Text())

	}

	return
}
