package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	lexer             Lexer
	defaultPrecedence int
	isOperator		  bool
	isBinaryOp		  bool
	binOpPrecedence   map[string]int
}

func NewParser(data string) Parser {
	return Parser{
		lexer: Lexer{
			Source:        data,
			CurrentChar:   -1,
			LastChar:      32,
			ignoreNewLine: true,
			ignoreSpace: true,
		},
		binOpPrecedence: map[string]int{
			"=":  2,
			"==": 9,
			"!=": 9,
			"<":  10,
			">":  10,
			">=": 10,
			"<=": 10,
			"+":  20,
			"-":  20,
			"*":  40,
			"/":  40,
		},
	}
}

func (p *Parser) getType(t string) string {
	switch t {
	case LitVoid:
		return LitVoid
	case LitFloat:
		return LitFloat
	case LitString:
		return LitString
	case LitBool:
		return LitBool
	case LitInt:
		return LitInt
	default:
		panic(fmt.Sprintf("type-%s-does-no-exit", t))
	}
}

func (p *Parser) parseArgs(callee bool) []ArgsPrototype {
	argsNames := []ArgsPrototype{}

	p.lexer.NextToken()
	for ; ; {
		if p.lexer.CurrentToken.kind == TokIdentifier {
			name := p.lexer.Identifier
			if callee {
				argsNames = append(argsNames, ArgsPrototype{
					Name:    name,
					ArgType: LitInt,
				})
			} else {
				p.lexer.NextToken()

				if p.lexer.CurrentToken.val != ':' {
					panic("wrong-args-definition")
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
				panic("wrong-args-definition")
			}
		} else {
			if !callee {
				panic("wrong-args-definition")
			}
		}

		p.lexer.NextToken()
		if p.lexer.CurrentToken.val != ',' {
			if p.lexer.CurrentToken.kind != TokRParen {
				panic("wrong-args-definition")
			}

			break
		}
		p.lexer.NextToken()
	}

	return argsNames
}

// TODO: remove callee bool
func (p *Parser) ParsePrototype(callee bool) (PrototypeAST, error) {
	p.lexer.ignoreAtoms = true

	pos := p.lexer.CurrentChar
	isOperator := p.isOperator
	isBinOp := p.isBinaryOp
	defPrecedence := p.defaultPrecedence

	p.isOperator = false
	p.isBinaryOp = false
	p.defaultPrecedence = 0

	if p.lexer.CurrentToken.kind != TokIdentifier {
		return PrototypeAST{}, errors.New("no-function-name")
	}

	funcName := p.lexer.Identifier
	p.lexer.NextToken()

	argsNames := []ArgsPrototype{}
	if p.lexer.CurrentToken.kind == TokLParen {
		argsNames = p.parseArgs(callee)
		if p.lexer.CurrentToken.kind != TokRParen {
			return PrototypeAST{}, errors.New(")-expected")
		}
		p.lexer.NextToken()
	}

	returnType := LitVoid
	if p.lexer.CurrentToken.val == ':' {
		p.lexer.NextToken()

		if p.lexer.CurrentToken.kind != TokIdentifier {
			return PrototypeAST{}, errors.New("return-type")
		}

		returnType = p.getType(p.lexer.Identifier)

		p.lexer.NextToken()
	}

	p.lexer.ignoreAtoms = false

	return PrototypeAST{
		position(pos),
		astPrototype,
		funcName,
		argsNames,
		isOperator,
		isBinOp,
		defPrecedence,
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

	if proto.ReturnType == LitVoid {
		body = append(body, &ReturnAST{
			position(pos),
			astReturn,
			nil,
		})
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
	return p.ParsePrototype(false)
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
		false,
		false,
		0,
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

func (p *Parser) checkBinOpPrec() (int, bool, string) {
	switch p.lexer.CurrentToken.kind {
	case TokAssign:
		prec, ok := p.binOpPrecedence["="]
		return prec, ok, "="
	case TokEqual:
		prec, ok := p.binOpPrecedence["=="]
		return prec, ok, "=="
	default:
		val := string(rune(p.lexer.CurrentToken.val))
		prec, ok := p.binOpPrecedence[val]
		return prec, ok, val
	}
}

func (p *Parser) ParseBinOpRHS(expressionPrec int, lhs AST) AST {
	pos := p.lexer.CurrentChar
	for ; ; {
		tokenPrec, ok, binop := p.checkBinOpPrec()

		if !ok {
			tokenPrec = -1
		}

		if tokenPrec < expressionPrec {
			return lhs
		}

		p.lexer.NextToken()

		rhs := p.ParsePrimary()

		if rhs == nil {
			return nil
		}

		nextPrec, ok, _ := p.checkBinOpPrec()

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
			binop,
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
	case TokReturn:
		return p.parseReturn()
	case TokLoop:
		return p.parseLoop()
	case TokBoolean:
		return p.parseBool()
	case TokIn:
		panic("Syntax Error: Invalid 'in' keyword usage")
	case TokFunction:
		panic("Syntax Error: Cannot define function here")
	case TokEnd:
		panic("Syntax Error: Extra end")
	case TokRParen:
		panic("Syntax Error: Invalid use of ')'")
	default:
		println("NUL", p.lexer.CurrentToken.kind, string(rune(p.lexer.CurrentToken.val)))
		p.lexer.NextToken()
		return nil
	}
}

func (p *Parser) parseBool() AST {
	pos := p.lexer.CurrentChar
	val := 0

	if p.lexer.Identifier == "true" {
		val = 1
	} else if p.lexer.Identifier != "false" {
		panic("Error occurred while parsing boolean")
	}

	p.lexer.NextToken()
	return &BoolAST{position(pos), astBool, val}
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
		if p.lexer.CurrentToken.kind == TokEOF {
			panic("Syntax Error: Function call is not closed")
		}

		arg := p.ParseExpression()
		if arg != nil {
			args = append(args, arg)
		}

		if p.lexer.CurrentToken.kind == TokRParen {
			break
		}

		if p.lexer.CurrentToken.val != ',' {
			panic("Syntax Error: Invalid character")
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
	kind := astNumberInt

	if p.lexer.IsFloat {
		kind = astNumberFloat
	}

	p.lexer.NextToken()

	return &NumberLiteralAST{position(pos), kind, val}
}

func (p *Parser) parseIfElse() AST {
	pos := p.lexer.CurrentChar
	p.lexer.NextToken()

	cond := p.ParseExpression()
	if cond == nil {
		panic("Syntax Error: No condition inside if")
	}

	trueBody := []AST{}
	for ; p.lexer.CurrentToken.kind != TokElse && p.lexer.CurrentToken.kind != TokElif && p.lexer.CurrentToken.kind != TokEnd; {
		if p.lexer.CurrentToken.kind == TokEOF {
			panic("Syntax Error: No end")
		}

		body := p.ParseExpression()
		if body != nil {
			trueBody = append(trueBody, body)
		}
	}

	elifBody := []ElifAST{}
	for ; ; {
		if p.lexer.CurrentToken.kind == TokEOF {
			panic("Syntax Error: No end")
		}

		if p.lexer.CurrentToken.kind != TokElif {
			break
		}

		p.lexer.NextToken()

		elifCond := p.ParseExpression()
		if elifCond == nil {
			panic("Syntax Error: No condition inside elif")
		}

		tempBody := []AST{}
		for ; p.lexer.CurrentToken.kind != TokEnd && p.lexer.CurrentToken.kind != TokElse && p.lexer.CurrentToken.kind != TokElif; {
			if p.lexer.CurrentToken.kind == TokEOF {
				panic("Syntax Error: No end")
			}

			body := p.ParseExpression()
			if body != nil {
				tempBody = append(tempBody, body)
			}
		}

		elifBody = append(elifBody, ElifAST{
			position(pos),
			astIfElse,
			elifCond,
			BlockAST{
				position(pos),
				astBlock,
				tempBody,
			},
		})
	}

	falseBody := []AST{}
	if p.lexer.CurrentToken.kind == TokElse {
		p.lexer.NextToken()
		for ; p.lexer.CurrentToken.kind != TokEnd; {
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
		elifBody,
	}
}

func (p *Parser) parseReturn() AST {
	pos := p.lexer.CurrentChar
	p.lexer.ignoreNewLine = false
	p.lexer.NextToken()
	p.lexer.ignoreNewLine = true

	if p.lexer.CurrentToken.val == 10 || p.lexer.CurrentToken.val == 13 {
		return &ReturnAST{
			position(pos),
			astReturn,
			nil,
		}
	}

	value := p.ParseExpression()
	return &ReturnAST{
		position(pos),
		astReturn,
		value,
	}
}

func (p *Parser) parseLoopBody() []AST {
	body := []AST{}
	for ; p.lexer.CurrentToken.kind != TokEnd; {
		if p.lexer.CurrentToken.kind == TokEOF {
			panic("Syntax Error: Loop has no end")
		}

		expr := p.ParseExpression()

		if expr != nil {
			body = append(body, expr)
		}
	}

	return body
}

func (p *Parser) parseLoop() AST {
	pos := p.lexer.CurrentChar
	p.lexer.NextToken()

	if p.lexer.CurrentToken.kind == TokBoolean {
		cond := p.ParseExpression()
		if cond == nil {
			panic("Syntax Error: No condition after 'in' keyword")
		}

		body := p.parseLoopBody()
		p.lexer.NextToken()

		return &LoopAST{
			position(pos),
			astLoop,
			false,
			cond,
			"",
			"",
			BlockAST{
				position(pos),
				astBlock,
				body,
			},
		}
	}

	if p.lexer.CurrentToken.kind != TokIdentifier {
		panic("Syntax Error: No index variable in the loop")
	}
	ind := p.lexer.Identifier

	p.lexer.NextToken()
	if p.lexer.CurrentToken.val != ',' {
		body := p.parseLoopBody()
		p.lexer.NextToken()

		return &LoopAST{
			position(pos),
			astLoop,
			false,
			&VariableAST{
				position: position(pos),
				kind:     astVariable,
				Name:     ind,
			},
			"",
			"",
			BlockAST{
				position(pos),
				astBlock,
				body,
			},
		}
	}

	p.lexer.NextToken()
	if p.lexer.CurrentToken.kind != TokIdentifier {
		panic("Syntax Error: No variable in the loop")
	}
	element := p.lexer.Identifier

	p.lexer.NextToken()
	if p.lexer.CurrentToken.kind != TokIn {
		panic("Syntax Error: No `in` keyword in the loop")
	}

	p.lexer.NextToken()
	cond := p.ParseExpression()
	if cond == nil {
		panic("Syntax Error: No condition after 'in' keyword")
	}

	body := p.parseLoopBody()
	p.lexer.NextToken()

	return &LoopAST{
		position(pos),
		astLoop,
		true,
		cond,
		ind,
		element,
		BlockAST{
			position(pos),
			astBlock,
			body,
		},
	}
}

func (p *Parser) parseAssign(panicMessage string) interface{} {
	p.lexer.NextToken()
	if p.lexer.CurrentToken.kind != TokAssign {
		panic(panicMessage)
	}

	p.lexer.NextToken()
	if p.lexer.CurrentToken.kind == TokIdentifier ||  p.lexer.CurrentToken.kind == TokAtom {
		return p.lexer.Identifier
	}

	if p.lexer.CurrentToken.kind == TokNumber {
		return p.lexer.numVal
	}

	panic(panicMessage)
}

func (p *Parser) parsePrimitiveAttr() {
	_ = p.lexer.CurrentChar
	p.lexer.NextToken()

	if p.lexer.CurrentToken.kind != TokLParen {
		panic("Syntax Error: No ( in the primitive attribute")
	}

	prevToken := 0
	for ;; {
		if p.lexer.CurrentToken.kind == TokRParen {
			break
		}
		p.lexer.NextToken()
		if p.lexer.CurrentToken.kind == TokRParen {
			if prevToken == ',' {
				panic("Syntax Error: Wrong attribute definition")
			}

			break
		}

		if p.lexer.CurrentToken.kind == TokEOF {
			panic("Syntax Error: Primitive attribute is not closed")
		}

		if p.lexer.CurrentToken.kind != TokIdentifier {
			panic("Syntax Error: No identifier in the primitive attribute")
		}

		switch p.lexer.Identifier {
		case "type":
			typ := p.parseAssign("Syntax Error: Invalid value assigning in the 'type' option of the primitive attribute")
			switch typ {
			case ":unary":
				p.isBinaryOp = false
			case ":binary":
				p.isBinaryOp = true
			default:
				panic("Syntax Error: Type '"+typ.(string)+"' in the primitive attribute does not exist")
			}
		case "precedence":
			precedence, ok := p.parseAssign("Syntax Error: Invalid value assigning in the 'precedence' option of the primitive attribute").(float64)
			if !ok {
				panic("Syntax Error: Could not assign value to precedence because value is not a number")
			}

			p.defaultPrecedence = int(precedence)
		default:
			panic("Syntax Error: There is no \"" + p.lexer.Identifier + "\" option in the primitive attribute")
		}

		p.lexer.NextToken()
		if p.lexer.CurrentToken.val != ',' && p.lexer.CurrentToken.kind != TokRParen {
			panic("Syntax Error: Wrong attribute definition")
		}

		prevToken = p.lexer.CurrentToken.val
	}

	p.isOperator = true
}

func (p *Parser) parseAttribute() {
	_ = p.lexer.CurrentChar

	p.lexer.NextToken()

	for ;; {
		if p.lexer.CurrentToken.val == ']' {
			break
		}

		if p.lexer.CurrentToken.kind == TokEOF {
			panic("Syntax Error: Attribute is not closed")
		}

		if p.lexer.CurrentToken.kind != TokIdentifier {
			panic("Syntax Error: No identifier in the attribute")
		}

		switch p.lexer.Identifier {
		case "primitive":
			p.parsePrimitiveAttr()
		default:
			panic("Attribute Error: '" + p.lexer.Identifier + "' does not exist")
		}

		p.lexer.NextToken()
		if p.lexer.CurrentToken.val != ',' && p.lexer.CurrentToken.val != ']' {
			panic("Syntax Error: Wrong attribute definition")
		}
	}

	p.lexer.NextToken()
}
