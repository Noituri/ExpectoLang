package main

import (
	"ExpectoLang/llvm/bindings/go/llvm"
)

const (
	astProcedure kind = iota
	astBinary
	astNumber
	astVariable
	astCall
	astBlock
	astPrototype
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
	Value  float64
}

type BinaryAST struct {
	position
	kind
	Op       rune
	Lhs, Rhs AST
}

type VariableAST struct {
	position
	kind
	Name string
}

type CallAST struct {
	position
	kind
	Callee string
	args []AST
}

type BlockAST struct {
	position
	kind
	Body []AST
}

type ArgsPrototype struct {
	Name string
	ArgType string
}

type PrototypeAST struct {
	position
	kind
	Name 	   string
	Args 	   []ArgsPrototype
	ReturnType string
}

type ProcedureAST struct {
	position
	kind
	Proto PrototypeAST
	Body  BlockAST
}