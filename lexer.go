package main

import (
	"errors"
	"strconv"
	"unicode"
)

const (
	TokEOF        = iota // End of string/file        0
	TokFunction          // function				  1
	TokEnd               // end of statements etc	  2
	TokIdentifier        // Identifier				  3
	TokReturn            // procedure return		  4
	TokExtern            // extern procedure		  5
	TokNumber            // number					  6
	TokStr               // string					  7
	TokBoolean           // boolean					  8
	TokLParen            // (						  9
	TokRParen            // )						  10
	TokIf                // If						  11
	TokElse              // If else					  12
	TokElif              // if else if (elif)		  13
	TokLoop              // for/while/foreach loop	  14
	TokIn                // loop - element in array	  15
	TokEqual             // ==						  16
	TokAssign            // =						  17
	TokAttribute         // #[attr1 = 0]			  18
	TokAtom              // :atom					  19
	TokUnknown           // Not specified type		  20
)

type tokenType uint8

type Token struct {
	kind tokenType
	val  int
}

type Lexer struct {
	Source        string
	CurrentToken  Token
	Identifier    string
	numVal        float64
	strVal        string
	CurrentChar   int
	LastChar      uint8
	isEOF         bool
	ignoreNewLine bool
	ignoreSpace   bool
	ignoreAtoms	  bool
	IsFloat       bool
}

func (l *Lexer) nextChar() error {
	if l.CurrentChar+1 < len(l.Source) {
		l.CurrentChar++
		l.LastChar = l.Source[l.CurrentChar]
		return nil
	}

	return errors.New("eof")
}

func (l *Lexer) removeSpace() (stopLexing bool) {
	for ((l.LastChar == 32 || l.LastChar == '\t') && l.ignoreSpace) || ((l.LastChar == 10 || l.LastChar == 13) && l.ignoreNewLine) {
		if l.nextChar() != nil {
			l.CurrentToken.kind = TokEOF
			l.CurrentToken.val = -1
			return true
		}
	}

	return false
}

func (l *Lexer) isAlphabetic() (stopLexing bool) {
	if unicode.IsLetter(rune(l.LastChar)) {
		l.Identifier = string(rune(l.LastChar))

		if l.nextChar() != nil {
			l.CurrentToken.kind = TokEOF
			l.CurrentToken.val = -1
			return true
		}

		for unicode.IsLetter(rune(l.LastChar)) || l.LastChar == '_' {
			l.Identifier += string(rune(l.LastChar))

			if l.nextChar() != nil {
				break
			}
		}

		if l.Identifier == "fc" {
			l.CurrentToken.kind = TokFunction
			l.CurrentToken.val = -1
			return true
		}

		if l.Identifier == "end" {
			l.CurrentToken.kind = TokEnd
			l.CurrentToken.val = -1
			return true
		}

		if l.Identifier == "return" {
			l.CurrentToken.kind = TokReturn
			l.CurrentToken.val = -1
			return true
		}

		if l.Identifier == "true" || l.Identifier == "false" {
			l.CurrentToken.kind = TokBoolean
			l.CurrentToken.val = -1
			return true
		}

		if l.Identifier == "if" {
			l.CurrentToken.kind = TokIf
			l.CurrentToken.val = -1
			return true
		}

		if l.Identifier == "else" {
			l.CurrentToken.kind = TokElse
			l.CurrentToken.val = -1
			return true
		}

		if l.Identifier == "elif" {
			l.CurrentToken.kind = TokElif
			l.CurrentToken.val = -1
			return true
		}

		if l.Identifier == "loop" {
			l.CurrentToken.kind = TokLoop
			l.CurrentToken.val = -1
			return true
		}

		if l.Identifier == "in" {
			l.CurrentToken.kind = TokIn
			l.CurrentToken.val = -1
			return true
		}

		l.CurrentToken.kind = TokIdentifier
		l.CurrentToken.val = -1
		return true
	}
	return false
}

func (l *Lexer) isDigit() (stopLexing bool) {
	if unicode.IsNumber(rune(l.LastChar)) || l.LastChar == '.' {
		tempStr := ""
		l.IsFloat = l.LastChar == '.'

		for ; ; {
			tempStr += string(rune(l.LastChar))
			if l.nextChar() != nil {
				l.CurrentToken.kind = TokEOF
				l.CurrentToken.val = -1
				return true
			}

			if l.LastChar == '.' {
				if l.IsFloat {
					panic("Invalid use of '.'")
				} else {
					l.IsFloat = true
				}
			}

			if !unicode.IsNumber(rune(l.LastChar)) && l.LastChar != 46 {
				break
			}
		}

		l.numVal, _ = strconv.ParseFloat(tempStr, 64)
		l.CurrentToken.kind = TokNumber
		l.CurrentToken.val = -1
		return true
	}

	return false
}

func (l *Lexer) isComment() (stopLexing bool) {
	if l.LastChar == '/' {
		l.isEOF = l.nextChar() != nil
		if l.LastChar == '/' {
			for ; ; {
				if l.nextChar() != nil {
					l.CurrentToken.kind = TokEOF
					l.CurrentToken.val = -1
					return true
				}

				if l.LastChar == 13 || l.LastChar == 10 {
					break
				}
			}

			l.NextToken()
			return true
		} else if l.LastChar == '*' {
			for ; ; {
				if l.nextChar() != nil {
					l.CurrentToken.kind = TokEOF
					l.CurrentToken.val = -1
					return true
				}

				if l.LastChar == '*' {
					if l.nextChar() != nil {
						l.CurrentToken.kind = TokEOF
						l.CurrentToken.val = -1
						return true
					}

					if l.LastChar == '/' {
						if l.nextChar() != nil {
							l.CurrentToken.kind = TokEOF
							l.CurrentToken.val = -1
							return true
						}
						break
					}
				}
			}

			l.NextToken()
			return true
		} else {
			l.LastChar = '/'
			return false
		}
	}
	return false
}

func (l *Lexer) isParen() (stopLexing bool) {
	if l.LastChar == '(' {
		l.isEOF = l.nextChar() != nil

		l.CurrentToken.kind = TokLParen
		l.CurrentToken.val = -1
		return true
	}

	if l.LastChar == ')' {
		l.isEOF = l.nextChar() != nil

		l.CurrentToken.kind = TokRParen
		l.CurrentToken.val = -1
		return true
	}

	return false
}

func (l *Lexer) isExtern() (stopLexing bool) {
	if l.LastChar == '@' {
		l.isEOF = l.nextChar() != nil

		l.CurrentToken.kind = TokExtern
		l.CurrentToken.val = -1
		return true
	}

	return false
}

func (l *Lexer) isStr() (stopLexing bool) {
	if l.LastChar == '"' {
		if l.nextChar() != nil {
			l.CurrentToken.kind = TokEOF
			l.CurrentToken.val = -1
			return true
		}

		l.strVal = ""
		for l.LastChar != '"' {
			l.strVal += string(rune(l.LastChar))

			if l.nextChar() != nil {
				l.isEOF = true
				break
			}
		}

		l.isEOF = l.nextChar() != nil

		l.CurrentToken.kind = TokStr
		l.CurrentToken.val = -1
		return true
	}

	return false
}

func (l *Lexer) isEqual() (stopLexing bool) {
	if l.LastChar == '=' {
		if l.nextChar() != nil {
			l.CurrentToken.kind = TokEOF
			l.CurrentToken.val = -1
			return true
		}

		if l.LastChar != '=' {
			l.isEOF = l.nextChar() != nil
			l.CurrentToken.kind = TokAssign
			l.CurrentToken.val = -1
			return true
		}

		l.isEOF = l.nextChar() != nil
		l.CurrentToken.kind = TokEqual
		l.CurrentToken.val = -1
		return true
	}

	return false
}

func (l *Lexer) isAttribute() (stopLexing bool) {
	if l.LastChar == '#' {
		if l.nextChar() != nil {
			l.CurrentToken.kind = TokEOF
			l.CurrentToken.val = -1
			return true
		}

		if l.LastChar == '[' {
			l.isEOF = l.nextChar() != nil
			l.CurrentToken.kind = TokAttribute
			l.CurrentToken.val = -1
			return true
		}
	}

	return false
}

func (l *Lexer) isAtom() (stopLexing bool) {
	if l.LastChar == ':' {
		oldToken := l.CurrentToken
		oldChar := l.CurrentChar

		l.ignoreSpace = false
		l.ignoreNewLine = false

		if l.nextChar() != nil {
			l.CurrentToken.kind = TokEOF
			l.CurrentToken.val = -1
			l.ignoreSpace = true
			l.ignoreNewLine = true
			return true
		}

		atom := ""
		for unicode.IsLetter(rune(l.LastChar)) || l.LastChar == '_' {
			atom += string(rune(l.LastChar))

			if l.nextChar() != nil {
				break
			}
		}

		l.ignoreSpace = true
		l.ignoreNewLine = true

		if atom == "" {
			l.LastChar = ':'
			l.CurrentToken = oldToken
			l.CurrentChar = oldChar
			return false
		}

		l.CurrentToken.kind = TokAtom
		l.CurrentToken.val = -1
		l.Identifier = ":" + atom
		return true
	}

	return false
}

func (l *Lexer) NextToken() {
	if l.isEOF {
		l.CurrentToken.kind = TokEOF
		l.CurrentToken.val = -1
		return
	}

	if l.removeSpace() {
		return
	}

	if l.isAttribute() {
		return
	}

	if l.isAlphabetic() {
		return
	}

	if l.isExtern() {
		return
	}

	if l.isStr() {
		return
	}

	if l.isDigit() {
		return
	}

	if l.isEqual() {
		return
	}

	if !l.ignoreAtoms {
		if l.isAtom() {
			return
		}
	}

	if l.isComment() {
		return
	}

	if l.isParen() {
		return
	}

	tempChar := l.LastChar

	l.isEOF = l.nextChar() != nil

	l.CurrentToken.kind = TokUnknown
	l.CurrentToken.val = int(tempChar)
}
