package main

const (
	TokenKindOpenParen  = "open-paren"
	TokenKindCloseParen = "close-paren"
	TokenKindPlus       = "plus"
	TokenKindMinus      = "minus"
	TokenKindMul        = "mul"
	TokenKindDiv        = "div"
	TokenKindNumber     = "number"
	TokenKindError      = "error"
)

type Token struct {
	Kind  string
	Value string
}

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input string
	pos   int
	start int
	state stateFn
	token Token
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) HasNext() bool {
	return l.pos < len(l.input)
}

func (l *Lexer) Next() Token {
	l.state = expr
	for l.state != nil {
		l.state = l.state(l)
	}
	return l.token
}

func (l *Lexer) nextChar() byte {
	if l.pos >= len(l.input) {
		l.pos++
		return 0
	}
	r := l.input[l.pos]
	l.pos++
	return r
}

func (l *Lexer) emit(kind string) {
	l.token = Token{Kind: kind, Value: l.input[l.start:l.pos]}
	l.start = l.pos
}

func expr(l *Lexer) stateFn {
	switch r := l.nextChar(); {
	case r == ' ':
		l.start = l.pos
		return expr
	case '0' <= r && r <= '9':
		l.pos--
		return number
	case r == '(':
		l.emit(TokenKindOpenParen)
	case r == ')':
		l.emit(TokenKindCloseParen)
	case r == '+':
		l.emit(TokenKindPlus)
	case r == '-':
		l.emit(TokenKindMinus)
	case r == '*':
		l.emit(TokenKindMul)
	case r == '/':
		l.emit(TokenKindDiv)
	default:
		l.pos = len(l.input)
		l.emit(TokenKindError)
	}
	return nil
}

func number(l *Lexer) stateFn {
	if r := l.nextChar(); '0' <= r && r <= '9' {
		return number
	}
	l.pos--
	l.emit(TokenKindNumber)
	return nil
}
