package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	lexer           Lexer
	binOpPrecedence map[string]int
}

func (p *Parser) getType(t string) string {
	switch t {
	case LitVoid:
		panic("NOT IMPLEMENTED")
	case LitFloat:
		return LitFloat
	case LitString:
		return LitString
	default:
		panic(fmt.Sprintf("type-%s-does-no-exit", t))
	}
}

func (p *Parser) ParsePrototype(callee bool) (PrototypeAST, error) {
	pos := p.lexer.CurrentChar

	if p.lexer.CurrentToken.kind != TokIdentifier {
		return PrototypeAST{}, errors.New("no-function-name")
	}

	funcName := p.lexer.Identifier
	p.lexer.NextToken()

	if p.lexer.CurrentToken.kind != TokLParen {
		return PrototypeAST{}, errors.New("(-expected")
	}

	argsNames := []ArgsPrototype{}

	p.lexer.NextToken()
	for ; ; {
		if p.lexer.CurrentToken.kind == TokIdentifier {
			name := p.lexer.Identifier
			if callee {
				argsNames = append(argsNames, ArgsPrototype{
					Name:    name,
					ArgType: LitString,
				})
			} else {
				p.lexer.NextToken()

				if p.lexer.CurrentToken.val != ':' {
					return PrototypeAST{}, errors.New("wrong-args-definition")
				}

				p.lexer.NextToken()

				argsNames = append(argsNames, ArgsPrototype{
					Name:    name,
					ArgType: p.getType(p.lexer.Identifier),
				})
			}
		} else if p.lexer.CurrentToken.kind == TokRParen {
			if len(argsNames) == 0 {
				break
			} else {
				return PrototypeAST{}, errors.New("wrong-args-definition")
			}
		} else {
			if !callee {
				return PrototypeAST{}, errors.New("wrong-args-definition")
			}
		}

		p.lexer.NextToken()
		if p.lexer.CurrentToken.val != ',' {
			if p.lexer.CurrentToken.kind != TokRParen {
				return PrototypeAST{}, errors.New("wrong-args-definition")
			}

			break
		}
		p.lexer.NextToken()
	}

	if p.lexer.CurrentToken.kind != TokRParen {
		return PrototypeAST{}, errors.New(")-expected")
	}

	p.lexer.NextToken()

	returnType := LitVoid

	if p.lexer.CurrentToken.val == ':' {
		p.lexer.NextToken()

		if p.lexer.CurrentToken.kind != TokIdentifier {
			return PrototypeAST{}, errors.New("return-type")
		}

		returnType = p.getType(p.lexer.Identifier)
		p.lexer.NextToken()
	}

	return PrototypeAST{
		position(pos),
		astPrototype,
		funcName,
		argsNames,
		returnType,
	}, nil
}

func (p *Parser) ParseFunction() (FunctionAST, error) {
	pos := p.lexer.CurrentChar
	p.lexer.NextToken()
	proto, err := p.ParsePrototype(false)

	if err != nil {
		return FunctionAST{}, err
	}

	body := []AST{}

	for ; p.lexer.CurrentToken.kind != TokEnd; {
		if p.lexer.CurrentToken.kind == TokEOF {
			return FunctionAST{}, errors.New("no-end")
		}

		expr := p.ParseExpression()

		if expr != nil {
			body = append(body, expr)
		}
	}

	p.lexer.NextToken()

	block := BlockAST{
		position(pos),
		astBlock,
		body,
	}
	return FunctionAST{
		position(pos),
		astFunction,
		proto,
		block,
	}, nil
}

func (p *Parser) ParseExtern() (PrototypeAST, error) {
	p.lexer.NextToken()
	return p.ParsePrototype(true)
}

func (p *Parser) ParseTopLevelExpr() (FunctionAST, error) {
	pos := p.lexer.CurrentChar

	expr := p.ParseExpression()
	if expr == nil {
		return FunctionAST{}, errors.New("no-expression")
	}

	proto := PrototypeAST{
		position(pos),
		astPrototype,
		"",
		nil,
		LitVoid,
	}

	block := BlockAST{
		position(pos),
		astBlock,
		[]AST{expr},
	}

	return FunctionAST{
		position(pos),
		astFunction,
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
	for ; ; {
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
			rhs = p.ParseBinOpRHS(tokenPrec+1, rhs)
			if rhs == nil {
				return nil
			}
		}

		lhs = &BinaryAST{
			position(pos),
			astBinary,
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
	case TokStr:
		return p.parseStr()
	case TokNumber:
		return p.parseNumber()
	case TokLParen:
		return p.parseParen()
	case TokIf:
		return p.parseIfElse()
	case TokFunction:
		panic("Syntax Error: Cannot define function here")
	case TokEnd:
		panic("Syntax Error: Extra end")
	case TokRParen:
		panic("Syntax Error: Invalid use of ')'")
	default:
		println(p.lexer.CurrentToken.kind)
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
		panic("Syntax Error: Parenthesis are not closed")
	}

	p.lexer.NextToken()

	return val
}

func (p *Parser) parseIdentifier() AST {
	pos := p.lexer.CurrentChar
	name := p.lexer.Identifier

	p.lexer.NextToken()

	if p.lexer.CurrentToken.kind != TokLParen {
		return &VariableAST{position(pos), astVariable, name}
	}

	p.lexer.NextToken()

	args := []AST{}

	for ; p.lexer.CurrentToken.kind != TokRParen; {
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
	return &CallAST{position(pos), astCall, name, args}
}

func (p *Parser) parseStr() AST {
	pos := p.lexer.CurrentChar
	val := p.lexer.strVal

	p.lexer.NextToken()
	return &StringAST{position(pos), astString, val}
}

func (p *Parser) parseNumber() AST {
	pos := p.lexer.CurrentChar
	val := p.lexer.numVal

	p.lexer.NextToken()
	return &NumberLiteralAST{position(pos), astNumber, val}
}

func (p *Parser) parseIfElse() AST {
	pos := p.lexer.CurrentChar
	p.lexer.NextToken()

	cond := p.ParseExpression()
	if cond == nil {
		panic("Syntax Error: No condition inside If")
	}

	trueBody := []AST{}

	for ;p.lexer.CurrentToken.kind != TokElse && p.lexer.CurrentToken.kind != TokEnd; {
		if p.lexer.CurrentToken.kind == TokEOF {
			panic("Syntax Error: No end")
		}

		body := p.ParseExpression()
		if body != nil {
			trueBody = append(trueBody, body)
		}
	}

	falseBody := []AST{}

	if p.lexer.CurrentToken.kind == TokElse {
		p.lexer.NextToken()
		for ;p.lexer.CurrentToken.kind != TokEnd; {
			if p.lexer.CurrentToken.kind == TokEOF {
				panic("Syntax Error: No end")
			}

			body := p.ParseExpression()
			if body != nil {
				falseBody = append(falseBody, body)
			}
		}
	}

	p.lexer.NextToken()

	return &IfElseAST{
		position(pos),
		astIfElse,
		cond,
		BlockAST{
			position(pos),
			astBlock,
			trueBody,
		},
		BlockAST{
			position(pos),
			astBlock,
			falseBody,
		},
	}
}
