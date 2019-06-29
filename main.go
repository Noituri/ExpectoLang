package main

import (
	"io/ioutil"
)

/*
static void HandleTopLevelExpression() {
  // Evaluate a top-level expression into an anonymous function.
  if (auto FnAST = ParseTopLevelExpr()) {
    if (auto *FnIR = FnAST->codegen()) {
      fprintf(stderr, "Read top-level expression:");
      FnIR->print(errs());
      fprintf(stderr, "\n");
    }
  } else {
    // Skip token for error recovery.
    getNextToken();
  }
}
 */

func handleProcedure(parser *Parser) {
	procAST, err := parser.ParseProcedure()

	if err != nil {
		println("Procedure Parse Error: ", err.Error())
		return
	}

	procIR := procAST.codegen()

	if procIR.IsNil() {
		println("Procedure CodeGen Error: Could not create IR")
		return
	}

	procIR.Dump()
}

func handleExtern(parser *Parser) {
	protoAST, err := parser.ParseExtern()

	if err != nil {
		println("Extern Parse Error: ", err.Error())
		return
	}

	externIR := protoAST.codegen()

	if externIR.IsNil() {
		println("Extern CodeGen Error: Could not create IR")
		return
	}

	externIR.Dump()
}

func handleTopLevelExpression(parser *Parser) {
	topAST, err := parser.ParseTopLevelExpr()

	if err != nil {
		println("Top Level Expression Parse Error: ", err.Error())
		return
	}

	topIR := topAST.codegen()

	if topIR.IsNil() {
		println("Top Level Expression CodeGen Error: Could not create IR")
		return
	}

	topIR.Dump()
}


func handle(parser Parser) {
	parser.lexer.NextToken()

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

	handle(parser)
}
