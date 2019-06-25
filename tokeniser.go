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
	Identifier  string
	floatVal    float64
	currentChar = -1
)

func nextChar(source string, char *uint8) error {
	if currentChar+ 1 < len(source) {
		currentChar++
		*char = source[currentChar]
		return nil
	}

	return errors.New("eof")
}

func GetToken(source string) uint8 {
	var lastChar uint8 = 32

	// is space
	for lastChar == 32 {
		if nextChar(source, &lastChar) != nil {
			currentChar = -1
			return TokEOF
		}
	}

	// is alphabetic
	if unicode.IsLetter(rune(lastChar)) {
		Identifier = string(rune(lastChar))

		if nextChar(source, &lastChar) != nil {
			currentChar = -1
			return TokEOF
		}

		for unicode.IsLetter(rune(lastChar)) {
			Identifier += string(rune(lastChar))
			if nextChar(source, &lastChar) != nil {
				break
			}
		}

		currentChar = -1

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
		currentChar = -1
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
		currentChar = -1
		return TokEOF
	}

	return tempChar
}
