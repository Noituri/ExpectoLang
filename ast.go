package main

import "errors"

type AST interface {
	Position() position
	Kind() kind
}

type position int
type kind int

func (p position) Position() position {
	return p
}

func (k kind) Kind() kind {
	return k
}

type FloatAST struct {
	position
	kind
	Value  float64
}

type BinaryAST struct {
	position
	kind
	Op       rune
	Lhs, Rhs AST
}

type PrototypeAST struct {
	position
	kind
	Name string
	Args []string
}

type FunctionAST struct {
	position
	kind
	Proto PrototypeAST
	Body  AST
}


func ParsePrototype() (PrototypeAST, error) {
	pos := CurrentChar

	if CurrentToken != TokIdentifier {
		println("LOL")
		return PrototypeAST{}, errors.New("no-function-name")
	}

	funcName := Identifier
	GetNextToken()

	if rune(CurrentToken) != '(' {
		return PrototypeAST{}, errors.New("(-expected")
	}

	argsNames := []string{}

	for GetNextToken(); CurrentToken == TokIdentifier; GetNextToken() {
		argsNames = append(argsNames, Identifier)
	}

	if rune(CurrentToken) != ')' {
		return PrototypeAST{}, errors.New(")-expected")
	}

	GetNextToken()

	return PrototypeAST{
		position(pos),
		kind(TokFunction),
		funcName,
		argsNames,
	}, nil
}

func ParseExtern() (PrototypeAST, error) {
	GetNextToken()
	return ParsePrototype()
}

func ParseTopLevelExpr() (FunctionAST, error) {
	pos := CurrentChar

	expr := ParseExpression()

	if expr == nil {
		return FunctionAST{}, errors.New("no-expression")
	}

	proto := PrototypeAST{
		position(pos),
		kind(TokFunction),
		"",
		nil,
	}

	return FunctionAST{
		position(pos),
		kind(TokFunction),
		proto,
		expr,
	}, nil
}

func ParseFunction() (FunctionAST, error) {
	pos := CurrentChar
	GetNextToken()
	proto, err := ParsePrototype()

	if err != nil {
		return FunctionAST{}, err
	}

	expr := ParseExpression()

	if expr == nil {
		return FunctionAST{}, errors.New("no-expression")
	}

	return FunctionAST{
		position(pos),
		kind(TokFunction),
		proto,
		expr,
	}, nil
}

func ParseExpression() AST {
	lhs := ParsePrimary()

	if lhs == nil {
		return nil
	}

	return ParseBinOpRHS(0, lhs)
}

func ParseBinOpRHS(expressionPrec int, lhs AST) AST {
	pos := CurrentChar
	for ;; {
		tokenPrec := GetTokenPrecedence()

		if tokenPrec < expressionPrec {
			return lhs
		}

		binop := CurrentToken
		GetNextToken()

		rhs := ParsePrimary()

		if rhs == nil {
			return nil
		}

		nextPrec := GetTokenPrecedence()

		if tokenPrec < nextPrec {
			rhs = ParseBinOpRHS(tokenPrec + 1, rhs)
			if rhs == nil {
				return nil
			}
		}

		return BinaryAST{
			position(pos),
			kind(-1),
			rune(binop),
			lhs,
			rhs,
		}
	}
}

func ParsePrimary() AST {
	switch CurrentToken {
	case TokFloat:
		return parseFloat()
	default:
		GetNextToken()
		return nil
	}
}

func parseFloat() AST {
	pos := CurrentChar
	val := floatVal

	GetNextToken()

	return &FloatAST{position(pos), TokFloat, val}
}