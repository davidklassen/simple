package main

type Op int

const (
	Noop = iota
	OpUnaryMinus
	OpBinaryMinus
	OpBinaryPlus
	OpDiv
	OpMul
)

type BaseNode struct {
	Op Op
}

func (n BaseNode) GetOp() Op {
	return n.Op
}

func (n BaseNode) Precedence() int {
	switch n.Op {
	case OpBinaryMinus, OpBinaryPlus:
		return 1
	case OpDiv, OpMul:
		return 2
	case OpUnaryMinus:
		return 3
	}
	return 0
}

type UnaryNode struct {
	BaseNode
	Value Node
}

func NewUnaryNode(op Op, value Node) UnaryNode {
	return UnaryNode{
		BaseNode: BaseNode{Op: op},
		Value:    value,
	}
}

type BinaryNode struct {
	BaseNode
	Left  Node
	Right Node
}

func NewBinaryNode(op Op, left, right Node) BinaryNode {
	return BinaryNode{
		BaseNode: BaseNode{Op: op},
		Left:     left,
		Right:    right,
	}
}

type ValueNode struct {
	BaseNode
	Value Token
}

func NewValueNode(value Token) ValueNode {
	return ValueNode{
		BaseNode: BaseNode{Op: Noop},
		Value:    value,
	}
}
