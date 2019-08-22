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
	astAttribute
)

type AST interface {
	Position() position
	Kind() kind
	codegen() llvm.Value
}

type position int
type kind int

func (p position) Position() position {
	return p
}

func (k kind) Kind() kind {
	return k
}

type NumberLiteralAST struct {
	position
	kind
	Value float64
}

type BinaryAST struct {
	position
	kind
	Op       string
	Lhs, Rhs AST
}

type BoolAST struct {
	position
	kind
	Value int
}

type StringAST struct {
	position
	kind
	Value string
}

type VariableAST struct {
	position
	kind
	Name 	string
	VarType string
	Mutable bool
}

type ElifAST struct {
	position
	kind
	Condition AST
	Body 	  BlockAST
}

type IfElseAST struct {
	position
	kind
	Condition AST
	TrueBody  BlockAST
	FalseBody BlockAST
	ElifBody  []ElifAST
}

type LoopAST struct {
	position
	kind
	forIn	   bool
	Condition  AST
	IndexVar   string
	ElementVar string
	Body       BlockAST
}

type CallAST struct {
	position
	kind
	Callee string
	args   []AST
}

type ReturnAST struct {
	 position
	 kind
	 Body AST
}

type BlockAST struct {
	position
	kind
	Elements []AST
}

type ArgsPrototype struct {
	Name    string
	ArgType string
}

type PrototypeAST struct {
	position
	kind
	Name       string
	Args       []ArgsPrototype
	IsOperator bool
	IsUnaryOp  bool
	Precedence int
	ReturnType string
}

type FunctionAST struct {
	position
	kind
	Proto PrototypeAST
	Body  BlockAST
}