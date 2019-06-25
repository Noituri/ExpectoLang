package main

import (
	"github.com/chzyer/readline"
	"io"
	"log"
	"strings"
)

var BinopPrecedence = make(map[rune]int)

func GetTokenPrecedence() int {
	p, ok := BinopPrecedence[rune(CurrentToken)]

	if !ok {
		return -1
	}

	return p
}
//static void MainLoop() {
//while (1) {
//fprintf(stderr, "ready> ");
//switch (CurTok) {
//case tok_eof:
//return;
//case ';': // ignore top-level semicolons.
//getNextToken();
//break;
//case tok_def:
//HandleDefinition();
//break;
//case tok_extern:
//HandleExtern();
//break;
//default:
//HandleTopLevelExpression();
//break;
//}
//}
//}
func main() {
	BinopPrecedence['<'] = 10
	BinopPrecedence['+'] = 20
	BinopPrecedence['-'] = 20
	BinopPrecedence['*'] = 40

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31mExpectoLang»\033[0m ",
		HistoryFile:     "./history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		line = strings.TrimSpace(line)
		Source = line
		GetNextToken()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		switch CurrentToken {
		case TokEOF:
			return
		case TokFunction: {
			fun, err := ParseFunction()

			if err != nil {
				println(err.Error())
				return
			}

			println("name: ", fun.Proto.Name, " args: ", fun.Proto.Args[0])
		}
		}
	}
}
