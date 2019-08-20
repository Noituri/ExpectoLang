package main

import (
	"ExpectoLang/llvm/bindings/go/llvm"
	"io/ioutil"
)

func handleFunction(parser *Parser, init bool) {
	functionAST, err := parser.ParseFunction()
	if err != nil {
		panic("Function Parse Error: " + err.Error())
	}

	if init {
		protoIR := functionAST.Proto.codegen()
		if protoIR.IsNil() {
			panic("Extern CodeGen Error: Could not create IR")
		}
	} else {
		fcIR := functionAST.codegen()
		if fcIR.IsNil() {
			panic("Function CodeGen Error: Could not create IR")
		}
	}
}

func handleExtern(parser *Parser, init bool) {
	protoAST, err := parser.ParseExtern()
	if err != nil {
		panic("Extern Parse Error: " + err.Error())
	}

	if init {
		externIR := protoAST.codegen()
		if externIR.IsNil() {
			panic("Extern CodeGen Error: Could not create IR")
		}
	}
}

func handleTopLevelExpression(parser *Parser, init bool) {
	topAST, err := parser.ParseTopLevelExpr()
	if err != nil {
		panic("Top Level Expression Parse Error: " + err.Error())
	}

	if !init {
		topIR := topAST.codegen()
		if topIR.IsNil() {
			panic("Top Level Expression CodeGen Error: Could not create IR")
		}
	}
}

func handle(parser Parser, init bool) {
	switch parser.lexer.CurrentToken.kind {
	case TokEOF:
		return
	case TokFunction:
		{
			handleFunction(&parser, init)
			//fun, err := parser.ParseFunction()
			//
			//if err != nil {
			//	println("Error: ", err.Error())
			//	return
			//}
			//
			//b := &bytes.Buffer{}
			//memviz.Map(b, &fun)
			//println(b.String())
		}
	case TokExtern:
		handleExtern(&parser, init)
	default:
		{
			handleTopLevelExpression(&parser, init)
			//fun, err := parser.ParseTopLevelExpr()
			//if err != nil {
			//	handle(parser)
			//	return
			//}
			//b := &bytes.Buffer{}
			//memviz.Map(b, &fun)
			//println(b.String())
		}
	}

	handle(parser, init)
}

func main() {
	data, err := ioutil.ReadFile("./example.exp")
	if err != nil {
		panic(err.Error())
	}

	parser := NewParser(string(data))
	parser.lexer.NextToken()
	InitModuleAndPassManager()

	handle(parser, true)
	handle(parser, false)
	if llvm.VerifyModule(module, llvm.PrintMessageAction) != nil {
		panic("Failed to verify module")
	}

	module.Dump()
}
