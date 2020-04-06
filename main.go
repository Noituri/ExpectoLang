package main

import (
	"ExpectoLang/llvm/bindings/go/llvm"
	"io/ioutil"
)

func handleFunction(parser *Parser, init bool) {
	functionAST := parser.parseFunction()
	if init {
		protoIR := functionAST.Proto.codegen()
		if protoIR.IsNil() {
			panic("Proto CodeGen Error: Could not create IR")
		}
	} else {
		fcIR := functionAST.codegen()
		if fcIR.IsNil() {
			panic("Function CodeGen Error: Could not create IR")
		}
	}
}

func handleExtern(parser *Parser, init bool) {
	//protoAST, err := parser.ParseExtern()
	//if err != nil {
	//	panic("Extern Parse Error: " + err.Error())
	//}
	//
	//if init {
	//	externIR := protoAST.codegen()
	//	if externIR.IsNil() {
	//		panic("Extern CodeGen Error: Could not create IR")
	//	}
	//}
}

func handleTopLevelExpression(parser *Parser, init bool) {
	//topAST, err := parser.ParseTopLevelExpr()
	//if err != nil {
	//	panic("Top Level Expression Parse Error: " + err.Error())
	//}
	//
	//if !init {
	//	topIR := topAST.codegen()
	//	if topIR.IsNil() {
	//		panic("Top Level Expression CodeGen Error: Could not create IR")
	//	}
	//}
}

func handle(parser Parser, init bool) {
	switch parser.lexer.token {
	case TokEOF:
		return
	case TokFunction:
		handleFunction(&parser, init)
	case TokExtern:
		handleExtern(&parser, init)
	case TokAttribute:
		parser.parseAttribute()
	default:
		handleTopLevelExpression(&parser, init)
	}

	handle(parser, init)
}

func initParser(data string) {
	parser := NewParser(data)

	parser.initialize = true
	handle(parser, true)
	parser.initialize = false

	parser.lexer = NewLexer(data)
	handle(parser, false)
}

func main() {
	InitModuleAndPassManager()

	data, err := ioutil.ReadFile("./test.nv")
	if err != nil {
		panic(err.Error())
	}

	initParser(string(data))

	if llvm.VerifyModule(module, llvm.PrintMessageAction) != nil {
		panic("Failed to verify module")
	}

	module.Dump()
}
