package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	lexer 			Lexer
	binOpPrecedence map[string] int
}

func (p *Parser) ParsePrototype() (PrototypeAST, error) {
	pos := p.lexer.CurrentChar

	if p.lexer.CurrentToken.kind != TokIdentifier {
		return PrototypeAST{}, errors.New("no-function-name")
	}

	funcName := p.lexer.Identifier
	p.lexer.NextToken()

	if p.lexer.CurrentToken.kind != TokLParen {
		return PrototypeAST{}, errors.New("(-expected")
	}

	argsNames := []string{}

	p.lexer.NextToken()
	for ;; {
		if p.lexer.CurrentToken.kind == TokIdentifier {
			argsNames = append(argsNames, p.lexer.Identifier)
		}

		p.lexer.NextToken()
		if p.lexer.CurrentToken.val != ',' && p.lexer.CurrentToken.kind != TokIdentifier {
			break
		}
	}

	if p.lexer.CurrentToken.kind != TokRParen {
		return PrototypeAST{}, errors.New(")-expected")
	}

	p.lexer.NextToken()

	returnType := LIT_VOID

	if p.lexer.CurrentToken.val == ':' {
		p.lexer.NextToken()

		if p.lexer.CurrentToken.kind != TokIdentifier {
			return PrototypeAST{}, errors.New("return-type")
		}

		switch p.lexer.Identifier {
		case LIT_VOID:
			panic("NOT IMPLEMENTED")
		case LIT_FLOAT:
			returnType = LIT_FLOAT
		case LIT_STRING:
			panic("NOT IMPLEMENTED")
		default:
			return PrototypeAST{}, errors.New(fmt.Sprintf("type-%s-does-no-exit", p.lexer.Identifier))
		}
	}

	return PrototypeAST{
		position(pos),
		kind(TokProcedure),
		funcName,
		argsNames,
		returnType,
	}, nil
}

func (p *Parser) ParseExtern() (PrototypeAST, error) {
	p.lexer.NextToken()
	return p.ParsePrototype()
}

func (p *Parser) ParseTopLevelExpr() (ProcedureAST, error) {
	pos := p.lexer.CurrentChar

	expr := p.ParseExpression()

	if expr == nil {
		return ProcedureAST{}, errors.New("no-expression")
	}

	proto := PrototypeAST{
		position(pos),
		kind(TokProcedure),
		"",
		nil,
		LIT_VOID,
	}


	block := BlockAST{
		position(pos),
		kind(TokProcedure),
		[]AST{expr},
	}

	return ProcedureAST{
		position(pos),
		kind(TokProcedure),
		proto,
		block,
	}, nil
}

func (p *Parser) ParseProcedure() (ProcedureAST, error) {
	pos := p.lexer.CurrentChar
	p.lexer.NextToken()
	proto, err := p.ParsePrototype()

	if err != nil {
		return ProcedureAST{}, err
	}

	body := []AST{}
	for ;; {
		expr := p.ParseExpression()
		if expr == nil {
			return ProcedureAST{}, errors.New("no-expression")
		}

		body = append(body, expr)
		if p.lexer.CurrentToken.kind == TokEnd {
			break
		}
	}

	block := BlockAST{
		position(pos),
		kind(TokProcedure),
		body,
	}
	return ProcedureAST{
		position(pos),
		kind(TokProcedure),
		proto,
		block,
	}, nil
}

func (p *Parser) ParseExpression() AST {
	lhs := p.ParsePrimary()

	if lhs == nil {
		return nil
	}

	return p.ParseBinOpRHS(0, lhs)
}

func (p *Parser) ParseBinOpRHS(expressionPrec int, lhs AST) AST {
	pos := p.lexer.CurrentChar
	for ;; {
		tokenPrec, ok := p.binOpPrecedence[string(rune(p.lexer.CurrentToken.val))]

		if !ok {
			tokenPrec = -1
		}

		if tokenPrec < expressionPrec {
			return lhs
		}

		binop := p.lexer.CurrentToken.val
		p.lexer.NextToken()

		rhs := p.ParsePrimary()

		if rhs == nil {
			return nil
		}

		nextPrec, ok := p.binOpPrecedence[string(rune(p.lexer.CurrentToken.val))]

		if !ok {
			tokenPrec = -1
		}

		if tokenPrec < nextPrec {
			rhs = p.ParseBinOpRHS(tokenPrec + 1, rhs)
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

func (p *Parser) ParsePrimary() AST {
	switch p.lexer.CurrentToken.kind {
	case TokIdentifier:
		return p.parseIdentifier()
	case TokNumber:
		return p.parseNumber()
	case TokLParen:
		return p.parseParen()
	case TokEnd:
		panic("Syntax Error: End")
	default:
		println("owo", p.lexer.CurrentToken.val)
		p.lexer.NextToken()
		return nil
	}
}

func (p *Parser) parseParen() AST {
	p.lexer.NextToken()
	val := p.ParseExpression()
	if val == nil {
		return nil
	}

	if p.lexer.CurrentToken.kind != TokRParen {
		return nil
	}

	p.lexer.NextToken()

	return val
}

func (p *Parser) parseIdentifier() AST {
	pos := p.lexer.CurrentChar
	name := p.lexer.Identifier

	p.lexer.NextToken()

	if p.lexer.CurrentToken.kind != TokLParen {
		return NumberLiteralAST{position(pos), TokNumber, p.lexer.numVal}
	}

	p.lexer.NextToken()

	args := []AST{}

		for ;p.lexer.CurrentToken.kind != TokRParen; {
			arg := p.ParseExpression()
			if arg == nil {
				return nil
			}

			args = append(args, arg)

			if p.lexer.CurrentToken.kind == TokRParen {
				break
			}

			if p.lexer.CurrentToken.val != ',' {
				return nil
			}

			p.lexer.NextToken()
		}


	p.lexer.NextToken()

	return CallAST{position(pos), TokProcedure, name, args}
}

func (p *Parser) parseNumber() AST {
	pos := p.lexer.CurrentChar
	val := p.lexer.numVal

	p.lexer.NextToken()
	return &NumberLiteralAST{position(pos), TokNumber, val}
}