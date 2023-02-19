package parser

type Stack []Token

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(attr Token) {
	*s = append(*s, attr)
}

func (s *Stack) Pop() (Token, bool) {
	if s.IsEmpty() {
		return ATTRIBUTE_UNKNOWN, false
	}

	attr := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return attr, true
}
