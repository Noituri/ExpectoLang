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
	println(tokens[parser.lexer.token])
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

//func initParser(data string) {
//	parser := NewParser(data)
//	for parser.lexer.token != TokEOF {
//		if parser.lexer.token == TokAttribute {
//			parser.parseAttribute()
//		}
//		if parser.lexer.
//		parser.lexer.nextToken()
//	}
//}

func main() {
	data, err := ioutil.ReadFile("./test.nv")
	if err != nil {
		panic(err.Error())
	}

	parser := NewParser(string(data))
	InitModuleAndPassManager()

	handle(parser, true)
	parser.lexer = NewLexer(string(data))
	handle(parser, false)
	if llvm.VerifyModule(module, llvm.PrintMessageAction) != nil {
		panic("Failed to verify module")
	}

	module.Dump()
}
