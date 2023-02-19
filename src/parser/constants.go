package parser

type Attr string

const (
	HTML  Attr = "HTML"
	DIV   Attr = "DIV"
	START Attr = "START"
	END   Attr = "END"
	CLOSE Attr = "CLOSE"
	PARAM Attr = "@PARAM"
)

var AttrDict map[Attr][]byte = map[Attr][]byte{
	HTML:  []byte("HTML>"),
	DIV:   []byte("DIV>"),
	START: []byte("<"),
	END:   []byte("</"),
	CLOSE: []byte(">"),
}
