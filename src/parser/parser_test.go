package parser_test

import (
	"bytes"
	"io"
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
}
