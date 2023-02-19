package parser_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"buns.gay/src/parser"
	"github.com/stretchr/testify/assert"
)

func TestNewCum(t *testing.T) {
	t.Run("when the input is basic", func(t *testing.T) {
		fin := io.NopCloser(strings.NewReader(`
		@CREATE HTML
			OPEN html CLOSE
		@ENDCREATE 

		HTML END
		`))
		fout := new(bytes.Buffer)

		p := parser.New()
		err := p.ParseCum(fin, fout)
		assert.NoError(t, err)

		gotHtml, err := io.ReadAll(fout)
		assert.NoError(t, err)

		assert.Equal(t, "<html ></HTML>", string(gotHtml))
	})

	t.Run("when END is used too much", func(t *testing.T) {
		fin := io.NopCloser(strings.NewReader("HTML END END"))
		fout := new(bytes.Buffer)

		p := parser.New()
		err := p.ParseCum(fin, fout)
		assert.Error(t, err)
	})

	t.Run("when a parameter is used", func(t *testing.T) {
		fin := io.NopCloser(strings.NewReader(`
		@CREATE HTML
			OPEN html CLOSE
		@ENDCREATE 

		@CREATE DIV
			OPEN div class=" @PARAM CLASS " CLOSE
		@ENDCREATE 
		
		HTML DIV @CLASS some-class END END`,
		))
		fout := new(bytes.Buffer)

		p := parser.New()
		err := p.ParseCum(fin, fout)
		assert.NoError(t, err)

		assert.Equal(t, &parser.AttributeTemplate{
			AttrName: "DIV",
			Params: map[parser.Token]bool{
				"@CLASS": true,
			},
			Template: []parser.Token{
				"OPEN", "div", "class=\"", "@PARAM",
				"CLASS", "\"", "CLOSE",
			},
		}, p.TM["DIV"])

		assert.Equal(t, map[parser.Token][]byte{
			"@CLASS": []byte("some-class"),
		}, p.VariableStore)

		gotHtml, err := io.ReadAll(fout)
		assert.NoError(t, err)

		assert.Equal(t, `<html ><div class=" some-class" ></DIV></HTML>`, string(gotHtml))
	})

	t.Run("when a parameter is not used", func(t *testing.T) {
		fin := io.NopCloser(strings.NewReader(`
		@CREATE HTML
			OPEN html CLOSE
		@ENDCREATE 

		@CREATE DIV
			OPEN div class=" @PARAM CLASS " CLOSE
		@ENDCREATE 
		
		HTML DIV END END`,
		))
		fout := new(bytes.Buffer)

		p := parser.New()
		err := p.ParseCum(fin, fout)
		assert.NoError(t, err)

		assert.Equal(t, &parser.AttributeTemplate{
			AttrName: "DIV",
			Params: map[parser.Token]bool{
				"@CLASS": true,
			},
			Template: []parser.Token{
				"OPEN", "div", "class=\"", "@PARAM",
				"CLASS", "\"", "CLOSE",
			},
		}, p.TM["DIV"])

		gotHtml, err := io.ReadAll(fout)
		assert.NoError(t, err)

		assert.Equal(t, `<html ><div class=" " ></DIV></HTML>`, string(gotHtml))
	})
	t.Run("hardmode", func(t *testing.T) {
		fin, err := os.Open("../fixtures/hardmode.cum")
		assert.NoError(t, err)
		fout := new(bytes.Buffer)

		p := parser.New()
		err = p.ParseCum(fin, fout)
		assert.NoError(t, err)

		assert.Equal(t, map[parser.Token]*parser.AttributeTemplate{
			"HTML": {
				AttrName: "HTML",
				Params:   map[parser.Token]bool{},
				Template: []parser.Token{
					"OPEN", "html", "CLOSE",
				},
			},
			"BODY": {
				AttrName: "BODY",
				Params:   map[parser.Token]bool{},
				Template: []parser.Token{
					"OPEN", "body", "CLOSE",
				},
			},
			"DIV": {
				AttrName: "DIV",
				Params: map[parser.Token]bool{
					"@CLASS": true,
				},
				Template: []parser.Token{
					"OPEN", "div", "class=\"", "@PARAM",
					"CLASS", "\"", "CLOSE",
				},
			},
			"BOILER": {
				AttrName: "BOILER",
				Params: map[parser.Token]bool{
					"@TITLE":  true,
					"@INSIDE": true,
				},
				Template: []parser.Token{
					"HTML", "BODY",
					"DIV", "@PARAM", "TITLE", "END",
					"DIV", "@CLASS", "body", "@PARAM", "INSIDE", "END",
					"END", "END",
				},
			}}, p.TM)

		assert.Equal(t, map[parser.Token][]byte{
			"@CLASS": []byte("poem-body"),
			"@TITLE": []byte("Rabbit Rabbit"),
		}, p.VariableStore)

		gotHtml, err := io.ReadAll(fout)
		assert.NoError(t, err)
		f, _ := os.Create("../fixtures/hardmode_test.html")
		f.Write(gotHtml)

		wantHTML, _ := os.ReadFile("../fixtures/hardmode.html")

		assert.Equal(t, strings.ReplaceAll(string(wantHTML), "\t", ""), string(gotHtml))

	})

}
