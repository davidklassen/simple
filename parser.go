package main

import "fmt"

type Node interface {
	GetOp() Op
	Precedence() int
}

var noop = BaseNode{}

func isValue(tok Token) bool {
	return tok.Kind == TokenKindNumber
}

func isUnary(tok Token) bool {
	return tok.Kind == TokenKindMinus
}

func isBinary(tok Token) bool {
	switch tok.Kind {
	case TokenKindMinus, TokenKindPlus, TokenKindDiv, TokenKindMul:
		return true
	}
	return false
}

func getUnaryOp(tok Token) Op {
	if tok.Kind == TokenKindMinus {
		return OpUnaryMinus
	}
	panic(fmt.Sprintf("not an unary operator: %s", tok.Value))
}

func getBinaryOp(tok Token) Op {
	switch tok.Kind {
	case TokenKindMinus:
		return OpBinaryMinus
	case TokenKindPlus:
		return OpBinaryPlus
	case TokenKindDiv:
		return OpDiv
	case TokenKindMul:
		return OpMul
	}
	panic(fmt.Sprintf("not a binary operator: %s", tok.Value))
}

type stack []Node

func (s *stack) top() Node {
	return (*s)[len(*s)-1]
}

func (s *stack) push(node Node) {
	*s = append(*s, node)
}

func (s *stack) pop() Node {
	var n Node
	n, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]
	return n
}

type Parser struct {
	lexer     *Lexer
	next      Token
	operators stack
	operands  stack
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) Parse() (ast Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			ast = nil
			err = fmt.Errorf("parse error: %s", r)
		}
	}()
	p.operators = []Node{noop}
	p.operands = make([]Node, 0)
	p.next = p.lexer.Next()
	p.expr()
	return p.operands.top(), nil
}

func (p *Parser) consume() {
	p.next = p.lexer.Next()
}

func (p *Parser) expect(kind string) {
	if p.next.Kind == kind {
		p.consume()
	} else {
		panic(fmt.Sprintf("unexpected token: %s", kind))
	}
}

func (p *Parser) expr() {
	p.partial()
	for isBinary(p.next) {
		p.pushOp(NewBinaryNode(getBinaryOp(p.next), nil, nil))
		p.consume()
		p.partial()
	}
	for p.operators.top() != noop {
		p.popOp()
	}
}

func (p *Parser) partial() {
	if isValue(p.next) {
		p.operands = append(p.operands, NewValueNode(p.next))
		p.consume()
	} else if p.next.Kind == TokenKindOpenParen {
		p.consume()
		p.operators = append(p.operators, noop)
		p.expr()
		p.expect(TokenKindCloseParen)
		p.operators.pop()
	} else if isUnary(p.next) {
		p.pushOp(NewUnaryNode(getUnaryOp(p.next), nil))
		p.consume()
		p.partial()
	} else {
		panic("invalid partial non-terminal")
	}
}

func (p *Parser) popOp() {
	if _, ok := p.operators.top().(BinaryNode); ok {
		op := p.operators.pop().(BinaryNode)
		op.Right = p.operands.pop()
		op.Left = p.operands.pop()
		p.operands.push(op)
	} else {
		op := p.operators.pop().(UnaryNode)
		op.Value = p.operands.pop()
		p.operands.push(op)
	}
}

func (p *Parser) pushOp(op Node) {
	for p.operators.top().Precedence() > op.Precedence() {
		p.popOp()
	}
	p.operators.push(op)
}
