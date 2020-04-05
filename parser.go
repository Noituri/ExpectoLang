package main

import (
	"fmt"
	"log"
)

var BinOpLookup = make(map[string][]ArgsPrototype)

type Parser struct {
	lexer             Lexer
	defaultPrecedence int
	isOperator        bool
	isBinaryOp        bool
	binOpPrecedence   map[string]int
	knownVars         map[string]string
}

func NewParser(data string) Parser {
	return Parser{
		lexer: NewLexer(data),
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

func (p *Parser) checkType(t string) string {
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

func (p *Parser) checkAndNext(tok Token) Pos {
	pos := p.lexer.pos
	if p.lexer.token != tok {
		log.Panicf("Invalid token. Expected: %s, got: %s", tokens[tok], tokens[p.lexer.token])
	}
	p.lexer.nextToken()
	return pos
}

func (p *Parser) parseArgs() []ArgsPrototype {
	var argsNames []ArgsPrototype

	for {
		if p.lexer.token == TokIdentifier {
			name := p.lexer.identifier
			p.lexer.nextToken()

			if p.lexer.token != TokTypeSpec {
				panic("After '"+name+"' argument there is no type specification.")
			}

			p.lexer.nextToken()
			_ = p.checkAndNext(TokIdentifier)
			varType := p.checkType(p.lexer.identifier)
			argsNames = append(argsNames, ArgsPrototype{
				Name:    name,
				ArgType: varType,
			})

			_, exist := p.knownVars[name]
			if exist {
				panic("Variable with the same name already exist.")
			}

			p.knownVars[name] = varType
		} else if p.lexer.token == TokRParen {
			if len(argsNames) == 0 {
				break
			} else {
				panic("Expected another argument.")
			}
		}

		if p.lexer.token != TokArgSep {
			if p.lexer.token != TokRParen {
				panic("Expected ')'")
			}

			break
		}
		p.lexer.nextToken()
	}

	return argsNames
}


func (p *Parser) parsePrototype() PrototypeAST {
	p.lexer.ignoreAtoms = true

	pos := p.lexer.pos
	isOperator := p.isOperator
	isBinOp := p.isBinaryOp
	defPrecedence := p.defaultPrecedence

	p.isOperator = false
	p.isBinaryOp = false
	p.defaultPrecedence = 0

	funcName := ""

	switch p.lexer.token {
	case TokIdentifier:
		if isOperator {
			panic("Error: Operator is not a special character")
		}

		funcName = p.lexer.identifier
		p.lexer.nextToken()
	default:
		if !isOperator {
			panic("Error: Only operators can use special character")
		}

		p.lexer.ignoreNewLine = false
		p.lexer.ignoreSpace = false
		for {
			if p.lexer.unknownVal == ' ' || p.lexer.unknownVal == '\n' || p.lexer.token == TokLParen {
				break
			}
			// ??? Not sure for what it was
			//if p.lexer.token != TokAssign && p.lexer.token != TokEqual {
			//	panic("Invalid operator name")
			//}

			if p.lexer.token == TokTypeSpec {
				panic("':' can't be used as operator name.")
			}

			if p.lexer.token == TokAssign {
				funcName += "="
			} else if p.lexer.token == TokEqual {
				funcName += "=="
			} else {
				funcName += string(p.lexer.unknownVal)
			}

			p.lexer.nextToken()
		}

		p.lexer.ignoreNewLine = true
		p.lexer.ignoreSpace = true

		if isBinOp {
			if _, found := p.binOpPrecedence[funcName]; found {
				// TODO: if binOp exists in current scope(!!!!) throw an exception
			}
			p.binOpPrecedence[funcName] = defPrecedence
		} else {
			if len(funcName) != 1 {
				panic("Error: Unary operator can only have one character name")
			}
		}

		if isBinOp {
			funcName = "binary_" + funcName
		} else {
			funcName = "unary_" + funcName
		}

		if p.lexer.token == TokUnknown {
			p.lexer.nextToken()
		}
	}

	var argsNames []ArgsPrototype
	if p.lexer.token == TokLParen {
		p.lexer.nextToken()
		argsNames = p.parseArgs()
		p.lexer.nextToken()
	}

	if isOperator && isBinOp && len(argsNames) != 2 {
		panic("Wrong number of arguments in the binary operator (" + funcName + ")")
	}

	if isOperator && !isBinOp && len(argsNames) != 1 {
		panic("Wrong number of arguments in the unary operator (" + funcName + ")")
	}

	if isBinOp {
		BinOpLookup[funcName] = argsNames
	}

	returnType := LitVoid
	if p.lexer.token == TokTypeSpec {
		p.lexer.nextToken()

		if p.lexer.token != TokIdentifier {
			panic("Expected a return type.")
		}

		returnType = p.checkType(p.lexer.identifier)
		p.lexer.nextToken()
	}

	p.lexer.ignoreAtoms = false

	return PrototypeAST{
		pos,
		astPrototype,
		funcName,
		argsNames,
		isOperator,
		isBinOp,
		defPrecedence,
		returnType,
	}
}

func (p *Parser) parseFunction() FunctionAST {
	pos := p.checkAndNext(TokFunction)
	p.knownVars = make(map[string]string)
	proto := p.parsePrototype()
	blockPos := p.checkAndNext(TokLBrace)

	var body []AST
	for p.lexer.token != TokRBrace {
		if p.lexer.token == TokEOF {
			panic("Function is not closed.")
		}
		expr := p.parseExpression()
		if expr != nil {
			body = append(body, expr)
		}
	}

	if proto.ReturnType == LitVoid {
		body = append(body, &ReturnAST{
			blockPos,
			astReturn,
			nil,
		})
	}
	p.lexer.nextToken()

	block := BlockAST{
		blockPos,
		astBlock,
		body,
	}
	return FunctionAST{
		pos,
		astFunction,
		proto,
		block,
	}
}

//func (p *Parser) parseExtern() PrototypeAST {
//	p.lexer.nextToken()
//	p.knownVars = make(map[string]string)
//	return p.parsePrototype()
//}
//
//func (p *Parser) parseTopLevelExpr() (FunctionAST, error) {
//	pos := p.lexer.CurrentChar
//
//	expr := p.ParseExpression()
//	if expr == nil {
//		return FunctionAST{}, errors.New("no-expression")
//	}
//
//	proto := PrototypeAST{
//		Pos(pos),
//		astPrototype,
//		"",
//		nil,
//		false,
//		false,
//		0,
//		LitVoid,
//	}
//
//	block := BlockAST{
//		Pos(pos),
//		astBlock,
//		[]AST{expr},
//	}
//
//	return FunctionAST{
//		Pos(pos),
//		astFunction,
//		proto,
//		block,
//	}, nil
//}

func (p *Parser) parseExpression() AST {
	lhs := p.parseUnary()
	if lhs == nil {
		return nil
	}

	return p.parseBinOpRHS(0, lhs)
}

func (p *Parser) checkBinOpPrec() (prec int, ok bool, operator string) {
	offset := p.lexer.offsetChar
	fwOffset := p.lexer.forwardOffset
	token := p.lexer.token
	lastChar := p.lexer.lastChar

	for {
		if p.lexer.token != TokAssign && p.lexer.token != TokEqual && p.lexer.token != TokUnknown {
			break
		}

		switch p.lexer.token {
		case TokAssign:
			operator += "="
		case TokEqual:
			operator += "=="
		default:
			operator += string(p.lexer.unknownVal)
		}

		//for k := range p.binOpPrecedence {
		//	if strings.HasPrefix(k, tempOp) {
		//		isMatch = true
		//		charCount = p.lexer.CurrentChar
		//		token = p.lexer.CurrentToken
		//		lastChar = p.lexer.LastChar
		//		break
		//	}
		//}
		//
		//if !isMatch {
		//	p.lexer.CurrentChar = charCount
		//	p.lexer.CurrentToken = token
		//	p.lexer.LastChar = lastChar
		//	break
		//}

		p.lexer.nextToken()
	}

	prec, ok = p.binOpPrecedence[operator]
	if !ok {
		p.lexer.offsetChar = offset
		p.lexer.forwardOffset = fwOffset
		p.lexer.token = token
		p.lexer.lastChar = lastChar
	}
	return prec, ok, operator
}

func (p *Parser) parseBinOpRHS(expressionPrec int, lhs AST) AST {
	pos := p.lexer.pos
	for {
		tokenPrec, ok, binop := p.checkBinOpPrec()
		if !ok {
			tokenPrec = -1
		}

		if tokenPrec < expressionPrec {
			return lhs
		}

		//p.lexer.nextToken()
		rhs := p.parseUnary()
		if rhs == nil {
			return nil
		}

		nextPrec, ok, _ := p.checkBinOpPrec()
		if !ok {
			tokenPrec = -1
		}

		if tokenPrec < nextPrec {
			rhs = p.parseBinOpRHS(tokenPrec+1, rhs)
			if rhs == nil {
				return nil
			}
		}

		lhs = &BinaryAST{
			Pos(pos),
			astBinary,
			binop,
			lhs,
			rhs,
		}
	}
}

func (p *Parser) parseStmt() AST {
	switch p.lexer.token {
	//case TokIdentifier:
	//	return p.parseIdentifier()
	//case TokStr:
	//	return p.parseStr()
	case TokNumber:
		return p.parseNumber()
	//case TokLParen:
	//	return p.parseParen()
	//case TokIf:
	//	return p.parseIfElse()
	//case TokReturn:
	//	return p.parseReturn()
	//case TokForLoop:
	//	return p.parseLoop()
	//case TokTrue, TokFalse:
	//	return p.parseBool()
	default:
		what := tokens[p.lexer.token]
		if p.lexer.token == TokUnknown {
			what = string(p.lexer.unknownVal)
		}
		panic("'"+what+"' has been used incorrectly.")
		return nil
	}
}

//func (p *Parser) parseBool() AST {
//	pos := p.lexer.CurrentChar
//	val := 0
//
//	if p.lexer.Identifier == "true" {
//		val = 1
//	} else if p.lexer.Identifier != "false" {
//		panic("Error occurred while parsing boolean")
//	}
//
//	p.lexer.nextToken()
//	return &BoolAST{Pos(pos), astBool, val}
//}
//
//func (p *Parser) parseParen() AST {
//	p.lexer.nextToken()
//	val := p.ParseExpression()
//	if val == nil {
//		return nil
//	}
//
//	if p.lexer.CurrentToken.kind != TokRParen {
//		panic("Syntax Error: Parenthesis are not closed")
//	}
//
//	p.lexer.nextToken()
//
//	return val
//}
//
//func (p *Parser) parseIdentifier() AST {
//	pos := p.lexer.CurrentChar
//	name := p.lexer.Identifier
//
//	p.lexer.nextToken()
//
//	if p.lexer.CurrentToken.kind != TokLParen {
//		varType, ok := p.knownVars[name]
//		if !ok {
//			panic(fmt.Sprintf(`Error: Variable "%s" does not exist!`, name))
//		}
//
//		return &VariableAST{
//			Pos: Pos(pos),
//			kind:     astVariable,
//			Name:     name,
//			VarType:  varType,
//		}
//	}
//
//	p.lexer.nextToken()
//
//	args := []AST{}
//
//	for ; p.lexer.CurrentToken.kind != TokRParen; {
//		if p.lexer.CurrentToken.kind == TokEOF {
//			panic("Syntax Error: Function call is not closed")
//		}
//
//		arg := p.ParseExpression()
//		if arg != nil {
//			args = append(args, arg)
//		}
//
//		if p.lexer.CurrentToken.kind == TokRParen {
//			break
//		}
//
//		if p.lexer.CurrentToken.val != ',' {
//			panic("Syntax Error: Invalid character")
//		}
//
//		p.lexer.nextToken()
//	}
//
//	p.lexer.nextToken()
//	return &CallAST{Pos(pos), astCall, name, args}
//}
//
//func (p *Parser) parseStr() AST {
//	pos := p.lexer.CurrentChar
//	val := p.lexer.strVal
//
//	p.lexer.nextToken()
//	return &StringAST{Pos(pos), astString, val}
//}

func (p *Parser) parseNumber() AST {
	pos := p.lexer.pos
	val := p.lexer.numVal
	kind := astNumberInt
	if p.lexer.isFloat {
		kind = astNumberFloat
	}

	p.lexer.nextToken()
	return &NumberLiteralAST{pos, kind, val}
}

//func (p *Parser) parseIfElse() AST {
//	pos := p.lexer.CurrentChar
//	p.lexer.nextToken()
//
//	cond := p.ParseExpression()
//	if cond == nil {
//		panic("Syntax Error: No condition inside if")
//	}
//
//	trueBody := []AST{}
//	for ; p.lexer.CurrentToken.kind != TokElse && p.lexer.CurrentToken.kind != TokElif && p.lexer.CurrentToken.kind != TokEnd; {
//		if p.lexer.CurrentToken.kind == TokEOF {
//			panic("Syntax Error: No end")
//		}
//
//		body := p.ParseExpression()
//		if body != nil {
//			trueBody = append(trueBody, body)
//		}
//	}
//
//	elifBody := []ElifAST{}
//	for ; ; {
//		if p.lexer.CurrentToken.kind == TokEOF {
//			panic("Syntax Error: No end")
//		}
//
//		if p.lexer.CurrentToken.kind != TokElif {
//			break
//		}
//
//		p.lexer.nextToken()
//
//		elifCond := p.ParseExpression()
//		if elifCond == nil {
//			panic("Syntax Error: No condition inside elif")
//		}
//
//		tempBody := []AST{}
//		for ; p.lexer.CurrentToken.kind != TokEnd && p.lexer.CurrentToken.kind != TokElse && p.lexer.CurrentToken.kind != TokElif; {
//			if p.lexer.CurrentToken.kind == TokEOF {
//				panic("Syntax Error: No end")
//			}
//
//			body := p.ParseExpression()
//			if body != nil {
//				tempBody = append(tempBody, body)
//			}
//		}
//
//		elifBody = append(elifBody, ElifAST{
//			Pos(pos),
//			astIfElse,
//			elifCond,
//			BlockAST{
//				Pos(pos),
//				astBlock,
//				tempBody,
//			},
//		})
//	}
//
//	falseBody := []AST{}
//	if p.lexer.CurrentToken.kind == TokElse {
//		p.lexer.nextToken()
//		for ; p.lexer.CurrentToken.kind != TokEnd; {
//			if p.lexer.CurrentToken.kind == TokEOF {
//				panic("Syntax Error: No end")
//			}
//
//			body := p.ParseExpression()
//			if body != nil {
//				falseBody = append(falseBody, body)
//			}
//		}
//	}
//
//	p.lexer.nextToken()
//
//	return &IfElseAST{
//		Pos(pos),
//		astIfElse,
//		cond,
//		BlockAST{
//			Pos(pos),
//			astBlock,
//			trueBody,
//		},
//		BlockAST{
//			Pos(pos),
//			astBlock,
//			falseBody,
//		},
//		elifBody,
//	}
//}
//
//func (p *Parser) parseReturn() AST {
//	pos := p.lexer.CurrentChar
//	p.lexer.ignoreNewLine = false
//	p.lexer.nextToken()
//	p.lexer.ignoreNewLine = true
//
//	if p.lexer.CurrentToken.val == 10 || p.lexer.CurrentToken.val == 13 {
//		return &ReturnAST{
//			Pos(pos),
//			astReturn,
//			nil,
//		}
//	}
//
//	value := p.ParseExpression()
//	return &ReturnAST{
//		Pos(pos),
//		astReturn,
//		value,
//	}
//}
//
//func (p *Parser) parseLoopBody() []AST {
//	body := []AST{}
//	for ; p.lexer.CurrentToken.kind != TokEnd; {
//		if p.lexer.CurrentToken.kind == TokEOF {
//			panic("Syntax Error: Loop has no end")
//		}
//
//		expr := p.ParseExpression()
//
//		if expr != nil {
//			body = append(body, expr)
//		}
//	}
//
//	return body
//}
//
//func (p *Parser) parseLoop() AST {
//	pos := p.lexer.CurrentChar
//	p.lexer.nextToken()
//
//	if p.lexer.CurrentToken.kind == TokBoolean {
//		cond := p.ParseExpression()
//		if cond == nil {
//			panic("Syntax Error: No condition after 'in' keyword")
//		}
//
//		body := p.parseLoopBody()
//		p.lexer.nextToken()
//
//		return &LoopAST{
//			Pos(pos),
//			astLoop,
//			false,
//			cond,
//			"",
//			"",
//			BlockAST{
//				Pos(pos),
//				astBlock,
//				body,
//			},
//		}
//	}
//
//	if p.lexer.CurrentToken.kind != TokIdentifier {
//		panic("Syntax Error: No index variable in the loop")
//	}
//	ind := p.lexer.Identifier
//
//	p.lexer.nextToken()
//	if p.lexer.CurrentToken.val != ',' {
//		body := p.parseLoopBody()
//		p.lexer.nextToken()
//
//		return &LoopAST{
//			Pos(pos),
//			astLoop,
//			false,
//			&VariableAST{
//				Pos: Pos(pos),
//				kind:     astVariable,
//				Name:     ind,
//				VarType:  LitBool,
//			},
//			"",
//			"",
//			BlockAST{
//				Pos(pos),
//				astBlock,
//				body,
//			},
//		}
//	}
//
//	p.lexer.nextToken()
//	if p.lexer.CurrentToken.kind != TokIdentifier {
//		panic("Syntax Error: No variable in the loop")
//	}
//	element := p.lexer.Identifier
//
//	p.lexer.nextToken()
//	if p.lexer.CurrentToken.kind != TokIn {
//		panic("Syntax Error: No `in` keyword in the loop")
//	}
//
//	p.lexer.nextToken()
//	cond := p.ParseExpression()
//	if cond == nil {
//		panic("Syntax Error: No condition after 'in' keyword")
//	}
//
//	// Shadowing variables
//	oldInd, okInd := p.knownVars[ind]
//	oldElement, okElem := p.knownVars[element]
//
//	p.knownVars[ind] = LitInt
//	p.knownVars[element] = ASTTypeToLit(cond.Kind())
//
//	body := p.parseLoopBody()
//
//	if okInd {
//		p.knownVars[ind] = oldInd
//	} else {
//		delete(p.knownVars, p.knownVars[ind])
//	}
//
//	if okElem {
//		p.knownVars[ind] = oldElement
//	} else {
//		delete(p.knownVars, p.knownVars[ind])
//	}
//
//	p.lexer.nextToken()
//
//	return &LoopAST{
//		Pos(pos),
//		astLoop,
//		true,
//		cond,
//		ind,
//		element,
//		BlockAST{
//			Pos(pos),
//			astBlock,
//			body,
//		},
//	}
//}
//
//func (p *Parser) parseAssign(panicMessage string) interface{} {
//	p.lexer.nextToken()
//	if p.lexer.CurrentToken.kind != TokAssign {
//		panic(panicMessage)
//	}
//
//	p.lexer.nextToken()
//	if p.lexer.CurrentToken.kind == TokIdentifier || p.lexer.CurrentToken.kind == TokAtom {
//		return p.lexer.Identifier
//	}
//
//	if p.lexer.CurrentToken.kind == TokNumber {
//		return p.lexer.numVal
//	}
//
//	panic(panicMessage)
//}
//
//func (p *Parser) parsePrimitiveAttr() {
//	_ = p.lexer.CurrentChar
//	p.lexer.nextToken()
//
//	if p.lexer.CurrentToken.kind != TokLParen {
//		panic("Syntax Error: No ( in the primitive attribute")
//	}
//
//	prevToken := 0
//	for ; ; {
//		if p.lexer.CurrentToken.kind == TokRParen {
//			break
//		}
//		p.lexer.nextToken()
//		if p.lexer.CurrentToken.kind == TokRParen {
//			if prevToken == ',' {
//				panic("Syntax Error: Wrong attribute definition")
//			}
//
//			break
//		}
//
//		if p.lexer.CurrentToken.kind == TokEOF {
//			panic("Syntax Error: Primitive attribute is not closed")
//		}
//
//		if p.lexer.CurrentToken.kind != TokIdentifier {
//			panic("Syntax Error: No identifier in the primitive attribute")
//		}
//
//		switch p.lexer.Identifier {
//		case "type":
//			typ := p.parseAssign("Syntax Error: Invalid value assigning in the 'type' option of the primitive attribute")
//			switch typ {
//			case ":unary":
//				p.isBinaryOp = false
//			case ":binary":
//				p.isBinaryOp = true
//			default:
//				panic("Syntax Error: Type '" + typ.(string) + "' in the primitive attribute does not exist")
//			}
//		case "precedence":
//			precedence, ok := p.parseAssign("Syntax Error: Invalid value assigning in the 'precedence' option of the primitive attribute").(float64)
//			if !ok {
//				panic("Syntax Error: Could not assign value to precedence because value is not a number")
//			}
//
//			p.defaultPrecedence = int(precedence)
//		default:
//			panic("Syntax Error: There is no \"" + p.lexer.Identifier + "\" option in the primitive attribute")
//		}
//
//		p.lexer.nextToken()
//		if p.lexer.CurrentToken.val != ',' && p.lexer.CurrentToken.kind != TokRParen {
//			panic("Syntax Error: Wrong attribute definition")
//		}
//
//		prevToken = p.lexer.CurrentToken.val
//	}
//
//	p.isOperator = true
//}
//
//func (p *Parser) parseAttribute() {
//	_ = p.lexer.CurrentChar
//
//	p.lexer.nextToken()
//
//	for ; ; {
//		if p.lexer.CurrentToken.val == ']' {
//			break
//		}
//
//		if p.lexer.CurrentToken.kind == TokEOF {
//			panic("Syntax Error: Attribute is not closed")
//		}
//
//		if p.lexer.CurrentToken.kind != TokIdentifier {
//			panic("Syntax Error: No identifier in the attribute")
//		}
//
//		switch p.lexer.Identifier {
//		case "primitive":
//			p.parsePrimitiveAttr()
//		default:
//			panic("Attribute Error: '" + p.lexer.Identifier + "' does not exist")
//		}
//
//		p.lexer.nextToken()
//		if p.lexer.CurrentToken.val != ',' && p.lexer.CurrentToken.val != ']' {
//			panic("Syntax Error: Wrong attribute definition")
//		}
//	}
//
//	p.lexer.nextToken()
//}

func (p *Parser) parseUnary() AST {
	pos := p.lexer.pos
	if p.lexer.token != TokUnknown || p.lexer.unknownVal == ' ' {
		return p.parseStmt()
	}

	unaryOp := p.lexer.unknownVal
	p.lexer.nextToken()
	if op := p.parseUnary(); op != nil {
		return &UnaryAST{
			Pos: pos,
			kind:     astUnary,
			Operator: int(unaryOp),
			Operand:  op,
		}
	}

	return nil
}