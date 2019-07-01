package main

import (
	"errors"
	"strconv"
	"unicode"
)

const (
	TokEOF      = iota // End of string/file
	TokFunction        // function
	TokEnd             // end of statements etc
	TokIdentifier
	TokReturn   // procedure return
	TokExtern   // extern procedure
	TokNumber   // number
	TokStr		// string
	TokLParen   // (
	TokRParen   // )
	TokIf		// If
	TokElse		// If else
	TokUnknown  // Not specified type
)

type tokenType uint8

type Token struct {
	kind tokenType
	val  int
}

type Lexer struct {
	Source       string
	CurrentToken Token
	Identifier   string
	numVal       float64
	strVal		 string
	CurrentChar  int
	LastChar     uint8
	isEOF        bool
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
	for l.LastChar == 32 || l.LastChar == 10 || l.LastChar == 13 || l.LastChar == '\t' {
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

		for unicode.IsLetter(rune(l.LastChar)) {
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

		l.CurrentToken.kind = TokIdentifier
		l.CurrentToken.val = -1
		return true
	}
	return false
}

func (l *Lexer) isDigit() (stopLexing bool) {
	if unicode.IsNumber(rune(l.LastChar)) || l.LastChar == '.' {
		tempStr := ""
		wasDot := l.LastChar == '.'

		for ;; {
			tempStr += string(rune(l.LastChar))
			if l.nextChar() != nil {
				l.CurrentToken.kind = TokEOF
				l.CurrentToken.val = -1
				return true
			}

			if l.LastChar == '.' {
				if wasDot {
					panic("Invalid use of '.'")
				} else {
					wasDot = true
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
		for  l.LastChar != '"' {
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

func (l *Lexer) NextToken() {
	if l.isEOF {
		l.CurrentToken.kind = TokEOF
		l.CurrentToken.val = -1
		return
	}

	if l.removeSpace() {
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
