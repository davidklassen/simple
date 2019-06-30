package main

import (
	"fmt"
	"strconv"
)

func Eval(ast Node) (int, error) {
	switch n := ast.(type) {
	case ValueNode:
		return evalValue(n)
	case UnaryNode:
		return evalUnary(n)
	case BinaryNode:
		return evalBinary(n)
	}
	return 0, nil
}

func evalBinary(n BinaryNode) (int, error) {
	var err error
	var left, right int
	left, err = Eval(n.Left)
	if err != nil {
		return 0, err
	}
	right, err = Eval(n.Right)
	if err != nil {
		return 0, err
	}
	switch n.Op {
	case OpBinaryMinus:
		return left - right, nil
	case OpBinaryPlus:
		return left + right, nil
	case OpDiv:
		return left / right, nil
	case OpMul:
		return left * right, nil
	}
	return 0, fmt.Errorf("invalid op: %d", n.Op)
}

func evalUnary(n UnaryNode) (int, error) {
	if n.Op == OpUnaryMinus {
		val, err := Eval(n.Value)
		if err != nil {
			return 0, err
		}
		return -val, nil
	}
	return 0, fmt.Errorf("invalid op: %d", n.Op)
}

func evalValue(n ValueNode) (int, error) {
	return strconv.Atoi(n.Value.Value)
}
