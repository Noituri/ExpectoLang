// GoLang Scanner inspired
package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token int

const (
	TokEOF        Token = iota // End of string/file
	TokIdentifier              // Identifier
	TokNumber                  // number
	TokStr                     // string
	TokLParen                  // (
	TokRParen                  // )
	TokLBrace                  // {
	TokRBrace                  // }
	TokEqual                   // ==
	TokAssign                  // =
	TokTypeSpec				   // : Used for specifying a type
	TokArgSep				   // , arg separator
	TokAttribute               // #[attr1 = 0]
	TokAtom                    // :atom

	KWBegin
	TokIn       // loop - element in array
	TokForLoop  // for/while/foreach loop
	TokElse     // If else
	TokIf       // If
	TokFalse    // boolean false
	TokTrue     // boolean true
	TokExtern   // extern procedure
	TokUnknown  // Not specified type
	TokReturn   // procedure return
	TokFunction // function
	KWEnd
)

var tokens = map[Token]string{
	TokAtom:       "ATOM",
	TokUnknown:    "UNKNOWN",
	TokEOF:        "EOF",
	TokIdentifier: "IDENT",
	TokNumber:     "NUMBER",
	TokStr:        "STRING",
	TokAttribute:  "ATTRIBUTE",
	TokExtern:     "@",
	TokFunction:   "fun",
	TokReturn:     "return",
	TokTrue:       "true",
	TokFalse:      "false",
	TokIf:         "if",
	TokElse:       "else",
	TokForLoop:    "for",
	TokIn:         "in",
	TokLParen:     "(",
	TokRParen:     ")",
	TokLBrace:     "{",
	TokRBrace:     "}",
	TokEqual:      "==",
	TokAssign:     "=",
	TokTypeSpec:   ":",
	TokArgSep:     ",",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := KWBegin + 1; i < KWEnd; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(identifier string) Token {
	if tok, found := keywords[identifier]; found {
		return tok
	}
	return TokIdentifier
}

type Lexer struct {
	source        string
	token         Token
	identifier    string
	unknownVal    rune
	numVal        float64
	strVal        string
	offsetChar    int
	forwardOffset int
	pos           Pos
	lastChar      rune
	isEOF         bool
	ignoreNewLine bool
	ignoreSpace   bool
	ignoreAtoms   bool
	isFloat       bool
}

func NewLexer(source string) Lexer {
	lexer := Lexer{
		source:        source,
		offsetChar:    0,
		forwardOffset: 0,
		pos:           Pos{col: 1, row: 0},
		ignoreNewLine: true,
		ignoreSpace:   true,
	}
	_ = lexer.nextChar()
	if lexer.lastChar == 0xFEFF {
		_ = lexer.nextChar()
	}
	lexer.nextToken()
	return lexer
}

func (l *Lexer) clone() Lexer {
	return *l
}

func (l *Lexer) nextChar() error {
	if l.forwardOffset < len(l.source) {
		l.offsetChar = l.forwardOffset
		if l.lastChar == '\n' {
			l.pos.col = 0
			l.pos.row++
		}
		l.pos.col++
		ch := rune(l.source[l.forwardOffset])
		addOffset := 1
		if ch == 0 {
			return errors.New("null character")
		} else if ch >= utf8.RuneSelf {
			ch, addOffset = utf8.DecodeRune([]byte(l.source[l.forwardOffset:]))
			if ch == utf8.RuneError && addOffset == 1 {
				return errors.New("illegal UTF-8 encoding")
			} else if ch == 0xFEFF && l.offsetChar > 0 {
				return errors.New("0xFEFF is only allowed as a first character")
			}
		}

		l.lastChar = ch
		l.forwardOffset += addOffset
		return nil
	}

	return errors.New("eof")
}

func (l *Lexer) peek() byte {
	if l.forwardOffset < len(l.source) {
		return l.source[l.forwardOffset]
	}
	return 0
}

func (l *Lexer) removeSpace() (stopLexing bool) {
	for ((l.lastChar == '\t' || l.lastChar == ' ') && l.ignoreSpace) || ((l.lastChar == 10 || l.lastChar == 13) && l.ignoreNewLine) {
		if l.nextChar() != nil {
			l.token = TokEOF
			return true
		}
	}

	return false
}

func (l *Lexer) isAlphabetic() (stopLexing bool) {
	if unicode.IsLetter(l.lastChar) {
		l.identifier = string(l.lastChar)

		if l.nextChar() != nil {
			return true
		}
		for unicode.IsLetter(l.lastChar) || l.lastChar == '_' {
			l.identifier += string(l.lastChar)
			if l.nextChar() != nil {
				break
			}
		}

		l.token = Lookup(l.identifier)
		return true
	}
	return false
}

func (l *Lexer) isDigit() (stopLexing bool) {
	if unicode.IsNumber(l.lastChar) || l.lastChar == '.' || l.lastChar == '-' {
		tempStr := ""
		l.isFloat = l.lastChar == '.'

		for {
			tempStr += string(l.lastChar)
			if l.nextChar() != nil {
				l.token = TokEOF
				return true
			}

			if l.lastChar == '.' {
				if l.isFloat {
					panic("Invalid use of '.'")
				} else {
					l.isFloat = true
				}
			}

			if !unicode.IsNumber(l.lastChar) && l.lastChar != 46 {
				break
			}
		}

		l.numVal, _ = strconv.ParseFloat(tempStr, 64)
		l.token = TokNumber
		return true
	}

	return false
}

func (l *Lexer) isComment() (stopLexing bool) {
	if l.lastChar == '/' {
		ch := l.peek()
		l.isEOF = ch == 0
		if ch == '/' {
			for {
				if l.nextChar() != nil {
					l.token = TokEOF
					return true
				}

				if l.lastChar == 13 || l.lastChar == 10 {
					break
				}
			}

			l.nextToken()
			return true
		} else if ch == '*' {
			for {
				if l.nextChar() != nil {
					l.token = TokEOF
					return true
				}

				if l.lastChar == '*' {
					if l.nextChar() != nil {
						l.token = TokEOF
						return true
					}

					if l.lastChar == '/' {
						if l.nextChar() != nil {
							l.token = TokEOF
							return true
						}
						break
					}
				}
			}

			l.nextToken()
			return true
		}
	}
	return false
}

func (l *Lexer) isParen() (stopLexing bool) {
	if l.lastChar == '(' {
		l.isEOF = l.nextChar() != nil
		l.token = TokLParen
		return true
	}

	if l.lastChar == ')' {
		l.isEOF = l.nextChar() != nil
		l.token = TokRParen
		return true
	}

	return false
}

func (l *Lexer) isBrace() (stopLexing bool) {
	if l.lastChar == '{' {
		l.isEOF = l.nextChar() != nil
		l.token = TokLBrace
		return true
	}

	if l.lastChar == '}' {
		l.isEOF = l.nextChar() != nil
		l.token = TokRBrace
		return true
	}

	return false
}

func (l *Lexer) isExtern() (stopLexing bool) {
	if l.lastChar == '@' {
		l.isEOF = l.nextChar() != nil
		l.token = TokExtern
		return true
	}

	return false
}

func (l *Lexer) isStr() (stopLexing bool) {
	if l.lastChar == '"' {
		if l.nextChar() != nil {
			l.token = TokEOF
			return true
		}

		l.strVal = ""
		for l.lastChar != '"' {
			l.strVal += string(l.lastChar)
			if l.nextChar() != nil {
				l.isEOF = true
				break
			}
		}

		l.strVal = strings.ReplaceAll(l.strVal, "\\r", "\r")
		l.strVal = strings.ReplaceAll(l.strVal, "\\n", "\n")
		l.strVal = strings.ReplaceAll(l.strVal, "\\b", "\b")
		l.strVal = strings.ReplaceAll(l.strVal, "\\a", "\a")
		l.strVal = strings.ReplaceAll(l.strVal, "\\f", "\f")
		l.strVal = strings.ReplaceAll(l.strVal, "\\t", "\t")
		l.strVal = strings.ReplaceAll(l.strVal, "\\v", "\v")
		l.strVal = strings.ReplaceAll(l.strVal, "\\\\", "\\")
		l.strVal = strings.ReplaceAll(l.strVal, "\\\"", "\"")

		l.isEOF = l.nextChar() != nil
		l.token = TokStr
		return true
	}

	return false
}

func (l *Lexer) isEqual() (stopLexing bool) {
	if l.lastChar == '=' {
		if l.nextChar() != nil {
			l.token = TokEOF
			return true
		}

		if l.lastChar != '=' {
			l.token = TokAssign
			return true
		}

		l.isEOF = l.nextChar() != nil
		l.token = TokEqual
		return true
	}

	return false
}

func (l *Lexer) isTypeSpec() (stopLexing bool) {
	if l.lastChar == ':' {
		l.isEOF = l.nextChar() != nil
		l.token = TokTypeSpec
		return true
	}

	return false
}

func (l *Lexer) isArgSep() (stopLexing bool) {
	if l.lastChar == ',' {
		l.isEOF = l.nextChar() != nil
		l.token = TokArgSep
		return true
	}

	return false
}

func (l *Lexer) isAttribute() (stopLexing bool) {
	if l.lastChar == '#' {
		if l.nextChar() != nil {
			l.token = TokEOF
			return true
		}

		if l.lastChar == '[' {
			l.isEOF = l.nextChar() != nil
			l.token = TokAttribute
			return true
		}
	}

	return false
}

func (l *Lexer) isAtom() (stopLexing bool) {
	if l.lastChar == ':' {
		oldToken := l.token
		oldOffset := l.offsetChar
		oldFwOffset := l.forwardOffset

		l.ignoreSpace = false
		l.ignoreNewLine = false

		if l.nextChar() != nil {
			l.token = TokEOF
			l.ignoreSpace = true
			l.ignoreNewLine = true
			return true
		}

		atom := ""
		for unicode.IsLetter(l.lastChar) || l.lastChar == '_' {
			atom += string(l.lastChar)

			if l.nextChar() != nil {
				break
			}
		}

		l.ignoreSpace = true
		l.ignoreNewLine = true

		if atom == "" {
			l.lastChar = ':'
			l.token = oldToken
			l.offsetChar = oldOffset
			l.forwardOffset = oldFwOffset
			return false
		}

		l.token = TokAtom
		l.identifier = ":" + atom
		return true
	}

	return false
}

func (l *Lexer) nextToken() {
	if l.isEOF {
		l.token = TokEOF
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

	if l.isBrace() {
		return
	}

	if l.isTypeSpec() {
		return
	}

	if l.isArgSep() {
		return
	}

	l.unknownVal = l.lastChar
	l.token = TokUnknown
	l.isEOF = l.nextChar() != nil
}
