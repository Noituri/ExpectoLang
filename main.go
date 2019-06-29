package main

import (
	"io/ioutil"
)

func handleProcedure(parser *Parser) {
	procAST, err := parser.ParseProcedure()

	if err != nil {
		panic("Procedure Parse Error: " + err.Error())
	}

	procIR := procAST.codegen()

	if procIR.IsNil() {
		panic("Procedure CodeGen Error: Could not create IR")
		return
	}
}

func handleExtern(parser *Parser) {
	protoAST, err := parser.ParseExtern()

	if err != nil {
		panic("Extern Parse Error: " + err.Error())
		return
	}

	externIR := protoAST.codegen()

	if externIR.IsNil() {
		panic("Extern CodeGen Error: Could not create IR")
		return
	}
}

func handleTopLevelExpression(parser *Parser) {
	//println(parser.lexer.CurrentToken.kind, parser.lexer.CurrentChar)
	topAST, err := parser.ParseTopLevelExpr()

	if err != nil {
		panic("Top Level Expression Parse Error: " + err.Error())
		return
	}

	topIR := topAST.codegen()

	if topIR.IsNil() {
		panic("Top Level Expression CodeGen Error: Could not create IR")
		return
	}
}


func handle(parser Parser) {
	switch parser.lexer.CurrentToken.kind {
	case TokEOF:
		return
	case TokProcedure: {
		handleProcedure(&parser)
		//fun, err := parser.ParseProcedure()
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
		handleExtern(&parser)
	default: {
		handleTopLevelExpression(&parser)
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

	handle(parser)
}

func main() {
	dat, err := ioutil.ReadFile("./example.exp")
	if err != nil {
		panic(err.Error())
	}

	parser := Parser{
		lexer: Lexer{
			Source: string(dat),
			CurrentChar: -1,
			LastChar: 32,
		},
		binOpPrecedence: map[string]int{
			"<": 10,
			"+": 20,
			"-": 20,
			"*": 40,
			"/": 40,
		},
	}

	parser.lexer.NextToken()
	handle(parser)
	module.Dump()
}
