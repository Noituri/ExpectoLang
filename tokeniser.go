package main

import (
	"strconv"
	"unicode"
)

const (
	tok_eof = iota
	tok_function
	tok_identifier
	tok_extern
	tok_float
)

var (
	Identifier  string
	floatVal    float64
	currectChar = 0
)

func GetToken(source string) uint8 {
	var lastChar uint8 = 32

	// is space
	for lastChar == 32 {
		lastChar = source[currectChar]
		currectChar++
	}

	// is alphabetic
	if unicode.IsLetter(rune(lastChar)) {
		Identifier = string(rune(lastChar))

		currectChar++
		lastChar = source[currectChar]
		for unicode.IsLetter(rune(lastChar)) {
			Identifier += string(rune(lastChar))

			currectChar++
			lastChar = source[currectChar]
		}

		if Identifier == "fun" {
			return tok_function
		}

		if Identifier == "@" {
			return tok_extern
		}

		return tok_identifier
	}

	// is digit or dot
	if unicode.IsNumber(rune(lastChar)) || lastChar == 46 {
		tempStr := ""

		//TODO check if there is second dot in number
		for ;; {
			tempStr += string(rune(lastChar))
			currectChar++
			lastChar = source[currectChar]

			if !unicode.IsNumber(rune(lastChar)) && lastChar != 46 {
				break
			}
		}

		floatVal, _ = strconv.ParseFloat(tempStr, 64)
		return tok_float
	}

	// Check for comments (# for now)
	if lastChar == 35 {
		for ;; {
			currectChar++
			lastChar = source[currectChar]

			if lastChar == 13 || lastChar == 10 {
				break
			}

			return GetToken(source)
		}
	}


	//TODO CHECK FOR EOF

	tempChar := lastChar
	currectChar++
	lastChar = source[currectChar]

	return tempChar
}
