package main

import (
	"errors"
	"strconv"
	"unicode"
)

const (
	TokEOF = iota
	TokFunction
	TokIdentifier
	TokExtern
	TokFloat
)

var (
	Source       string
	CurrentToken uint8
	Identifier   string
	floatVal     float64
	CurrentChar  = -1
)

func GetNextToken() {
	CurrentToken = GetToken(Source)
}

func nextChar(source string, char *uint8) error {
	if CurrentChar + 1 < len(source) {
		CurrentChar++
		*char = source[CurrentChar]
		return nil
	}

	return errors.New("eof")
}

var lastChar uint8 = 32

func GetToken(source string) uint8 {
	// is space
	for lastChar == 32 {
		if nextChar(source, &lastChar) != nil {
			CurrentChar = -1
			return TokEOF
		}
	}

	// is alphabetic
	if unicode.IsLetter(rune(lastChar)) {
		Identifier = string(rune(lastChar))

		if nextChar(source, &lastChar) != nil {
			CurrentChar = -1
			return TokEOF
		}

		for unicode.IsLetter(rune(lastChar)) {
			Identifier += string(rune(lastChar))

			if nextChar(source, &lastChar) != nil {
				break
			}
		}

		if Identifier == "fun" {
			return TokFunction
		}

		if Identifier == "extern" {
			return TokExtern
		}

		return TokIdentifier
	}

	// is digit or dot
	if unicode.IsNumber(rune(lastChar)) || lastChar == 46 {
		tempStr := ""

		//TODO check if there is second dot in number
		for ;; {
			tempStr += string(rune(lastChar))

			if !unicode.IsNumber(rune(lastChar)) && lastChar != 46 || nextChar(source, &lastChar) != nil {
				break
			}
		}

		floatVal, _ = strconv.ParseFloat(tempStr, 64)
		return TokFloat
	}

	eof := false

	// Check for comments (# for now)
	if lastChar == 35 {

		for ;; {
			if nextChar(source, &lastChar) != nil {
				eof = true
				break
			}

			if lastChar == 13 || lastChar == 10 {
				break
			}
		}

		if !eof {
			return GetToken(source)
		}
	}

	tempChar := lastChar

	if nextChar(source, &lastChar) != nil || eof {
		CurrentChar = -1
	}

	return tempChar
}
