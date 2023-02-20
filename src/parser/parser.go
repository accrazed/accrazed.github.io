package parser

import (
	"bufio"
	"fmt"
	"io"
)

type Parser struct {
	TM            map[Token]*AttributeTemplate
	VariableStore map[Token][]Token
	TokenStack    *Stack
}

type AttributeTemplate struct {
	AttrName Token
	Params   map[Token]bool
	Template []Token
}

func New() *Parser {
	return &Parser{
		TM:            map[Token]*AttributeTemplate{},
		VariableStore: map[Token][]Token{},
		TokenStack:    &Stack{},
	}
}

func (p *Parser) ParseCum(in io.Reader, w io.Writer) error {
	scn := bufio.NewScanner(in)
	scn.Split(bufio.ScanWords)

	for scn.Scan() {
		token := Token(scn.Text())

		if token == CREATE {
			p.ParseTemplate(scn)
			continue
		}

		if err := p.ParseTokens(w, scn, token); err != nil {
			return fmt.Errorf("couldn't parse attribute %s: %w", token, err)
		}
	}

	return nil
}

func (p *Parser) ParseTokens(w io.Writer, scn *bufio.Scanner, tokens ...Token) error {
TokenLoop:
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token {
		case OPEN:
			w.Write([]byte("<"))
		case CLOSE:
			w.Write([]byte(">"))
		case END:
			lastToken, ok := p.TokenStack.Pop()
			if !ok {
				return fmt.Errorf("unexpected END in stack: %v", p.TokenStack)
			}
			w.Write([]byte("</"))
			w.Write([]byte(p.TM[lastToken].AttrName))
			w.Write([]byte(">"))
		case PARAM:
			paramField := Token("@" + tokens[i+1])
			paramVals, ok := p.VariableStore[paramField]
			if !ok {
				paramVals = []Token{}
			}
			p.ParseTokens(w, scn, paramVals...)
			i++
		case TEXT:
			w.Write([]byte("<div>"))
		WordLoop:
			for scn.Scan() {
				word := scn.Text()

				switch word {
				case string(ENDTEXT):
					break WordLoop
				case `\`:
					w.Write([]byte("</div>\n"))
				case `/`:
					w.Write([]byte("<div>"))
				case string(BR):
					w.Write([]byte("<br/>\n"))
				default:
					w.Write([]byte(word + " "))
				}

			}
			w.Write([]byte("</div>"))

		default:
			// token is just a word
			attrTemplate, ok := p.TM[token]
			if !ok {
				w.Write([]byte(token))
				w.Write([]byte(" "))
				continue
			}

			// token is a template
			p.TokenStack.Push(token)
			for scn.Scan() {
				// check for params
				paramField := Token(scn.Text())
				_, ok := p.TM[token].Params[paramField]
				if !ok {
					p.ParseTokens(w, scn, append(attrTemplate.Template, paramField)...)
					continue TokenLoop
				}

				// param exists
				p.VariableStore[paramField] = []Token{}
				for scn.Scan() {
					paramValue := Token(scn.Text())
					fmt.Println(paramValue)

					if paramValue == ENDPARAM {
						break
					}
					p.VariableStore[paramField] =
						append(p.VariableStore[paramField], paramValue)

				}
			}

			p.ParseTokens(w, scn, attrTemplate.Template...)
		}
	}

	return nil
}

func (p *Parser) ParseTemplate(scn *bufio.Scanner) error {
	scn.Scan()
	varName := Token(scn.Text())
	scn.Scan()
	attrName := Token(scn.Text())
	template := &AttributeTemplate{
		AttrName: attrName,
		Template: []Token{},
		Params:   map[Token]bool{},
	}

	for scn.Scan() {
		attr := Token(scn.Text())

		if attr == ENDCREATE {
			break
		}

		template.Template = append(template.Template, attr)

		// log param field name
		if attr == PARAM {
			scn.Scan()
			attr = Token(scn.Text())
			template.Params["@"+attr] = true
			template.Template = append(template.Template, attr)
		}
	}

	p.TM[varName] = template

	return nil
}
