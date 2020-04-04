package main

import (
	"ExpectoLang/llvm/bindings/go/llvm"
)

const (
	astFunction kind = iota
	astBinary
	astNumberFloat
	astNumberInt
	astString
	astBool
	astVariable
	astCall
	astBlock
	astPrototype
	astIfElse
	astReturn
	astLoop
	astUnary
)

type AST interface {
	Position() Pos
	Kind() kind
	codegen() llvm.Value
}

type Pos struct {
	row int
	col int
}
type kind int

func (p Pos) Position() Pos {
	return p
}

func (k kind) Kind() kind {
	return k
}

type NumberLiteralAST struct {
	Pos
	kind
	Value float64
}

type BinaryAST struct {
	Pos
	kind
	Op       string
	Lhs, Rhs AST
}

type UnaryAST struct {
	Pos
	kind
	Operator int
	Operand  AST
}

type BoolAST struct {
	Pos
	kind
	Value int
}

type StringAST struct {
	Pos
	kind
	Value string
}

type VariableAST struct {
	Pos
	kind
	Name    string
	VarType string
	Mutable bool
}

type ElifAST struct {
	Pos
	kind
	Condition AST
	Body      BlockAST
}

type IfElseAST struct {
	Pos
	kind
	Condition AST
	TrueBody  BlockAST
	FalseBody BlockAST
	ElifBody  []ElifAST
}

type LoopAST struct {
	Pos
	kind
	forIn      bool
	Condition  AST
	IndexVar   string
	ElementVar string
	Body       BlockAST
}

type CallAST struct {
	Pos
	kind
	Callee string
	args   []AST
}

type ReturnAST struct {
	Pos
	kind
	Body AST
}

type BlockAST struct {
	Pos
	kind
	Elements []AST
}

type ArgsPrototype struct {
	Name    string
	ArgType string
}

type PrototypeAST struct {
	Pos
	kind
	Name       string
	Args       []ArgsPrototype
	IsOperator bool
	IsUnaryOp  bool
	Precedence int
	ReturnType string
}

type FunctionAST struct {
	Pos
	kind
	Proto PrototypeAST
	Body  BlockAST
}
