package main

import (
	"bytes"
	"github.com/bradleyjkemp/memviz"
	"io/ioutil"
)

func handle(parser Parser) {
	parser.lexer.NextToken()

	switch parser.lexer.CurrentToken.kind {
	case TokEOF:
		return
	case TokProcedure: {
		fun, err := parser.ParseProcedure()

		if err != nil {
			println("Error: ", err.Error())
			return
		}

		b := &bytes.Buffer{}
		memviz.Map(b, &fun)
		println(b.String())
	}
	default: {
		fun, err := parser.ParseTopLevelExpr()
		if err != nil {
			handle(parser)
			return
		}
		b := &bytes.Buffer{}
		memviz.Map(b, &fun)
		println(b.String())
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
