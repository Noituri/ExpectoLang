package main

import (
	"ExpectoLang/llvm/bindings/go/llvm"
)

const (
	astFunction kind = iota
	astBinary
	astNumber
	astString
	astVariable
	astCall
	astBlock
	astPrototype
	astIfElse
	astReturn
	astLoop
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

type StringAST struct {
	position
	kind
	Value string
}

type VariableAST struct {
	position
	kind
	Name string
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
	Condition AST
	Body 	  BlockAST
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
	ReturnType string
}

type FunctionAST struct {
	position
	kind
	Proto PrototypeAST
	Body  BlockAST
}
