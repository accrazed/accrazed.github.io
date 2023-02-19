package parser

type Token string

var (
	ATTRIBUTE_UNKNOWN Token = "ATTRIBUTE_UNKNOWN"

	OPEN  Token = "OPEN"
	END   Token = "END"
	CLOSE Token = "CLOSE"

	PARAM     Token = "@PARAM"
	TEXT      Token = "@TEXT"
	ENDTEXT   Token = "@ENDTEXT"
	CREATE    Token = "@CREATE"
	ENDCREATE Token = "@ENDCREATE"
)
