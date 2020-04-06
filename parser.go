package main

import (
	"fmt"
	"strings"
)

var BinOpLookup = make(map[string][]ArgsPrototype)

type Parser struct {
	lexer             Lexer
	defaultPrecedence int
	isOperator        bool
	isBinaryOp        bool
	binOpPrecedence   map[string]int
	knownVars         map[string]string
	initialize		  bool
	errors			  []string
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

func (p *Parser) addError(err string) {
	if !p.initialize {
		panic(err)
	}
	p.errors = append(p.errors, err)
	p.lexer.nextToken()
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
		p.addError(fmt.Sprintf("Type %s doesn't exist.", t))
		return ""
	}
}

func (p *Parser) checkAndNext(tok Token) Pos {
	pos := p.lexer.pos
	if p.lexer.token != tok {
		got := tokens[p.lexer.token]
		if p.lexer.token == TokUnknown {
			got = string(p.lexer.unknownVal)
		}
		p.addError(fmt.Sprintf("Invalid token. Expected: %s, got: %s", tokens[tok], got))
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
				p.addError("After '"+name+"' argument there is no type specification.")
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
				p.addError("Variable with the same name already exist.")
			}

			p.knownVars[name] = varType
		} else if p.lexer.token == TokRParen {
			if len(argsNames) == 0 {
				break
			} else {
				p.addError("Expected another argument.")
			}
		}

		if p.lexer.token != TokArgSep {
			if p.lexer.token != TokRParen {
				p.addError("Expected ')'")
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
			p.addError("Error: Operator is not a special character")
		}

		funcName = p.lexer.identifier
		p.lexer.nextToken()
	default:
		if !isOperator {
			p.addError("Only operators can use special character")
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
				p.addError("':' can't be used as operator name.")
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
				p.addError("Error: Unary operator can only have one character name")
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
		p.addError("Wrong number of arguments in the binary operator (" + funcName + ")")
	}

	if isOperator && !isBinOp && len(argsNames) != 1 {
		p.addError("Wrong number of arguments in the unary operator (" + funcName + ")")
	}

	if isBinOp {
		BinOpLookup[funcName] = argsNames
	}

	returnType := LitVoid
	if p.lexer.token == TokTypeSpec {
		p.lexer.nextToken()

		if p.lexer.token != TokIdentifier {
			p.addError("Expected a return type.")
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
			p.addError("Function is not closed.")
		}

		stmt := p.parseStmt()
		if stmt != nil {
			body = append(body, stmt)
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

// This solution is super not clean IMHO. I could not think of a better one tho
func (p *Parser) checkBinOpPrec() (prec, tokenOffset int, ok bool, operator string) {
	oldLexer := p.lexer.clone()
	var prevLexers []Lexer
	goBack := 0

	for p.lexer.token == TokAssign || p.lexer.token == TokEqual || p.lexer.token == TokUnknown {
		tempOperator := ""
		switch p.lexer.token {
		case TokAssign:
			tempOperator += "="
		case TokEqual:
			tempOperator += "=="
		default:
			tempOperator += string(p.lexer.unknownVal)
		}

		matched := false
		for op := range p.binOpPrecedence {
			if strings.HasPrefix(op, operator + tempOperator) {
				goBack++
				prevLexers = append(prevLexers, p.lexer.clone())
				operator += tempOperator
				p.lexer.nextToken()
				matched = true
				break
			}
		}

		if !matched {
			break
		}
	}

	prec, ok = p.binOpPrecedence[operator]
	if !ok {
		p.lexer = oldLexer
	} else {
		p.lexer = prevLexers[len(prevLexers) - goBack]
	}

	return prec, goBack, ok, operator
}

func (p *Parser) parseBinOpRHS(expressionPrec int, lhs AST) AST {
	pos := p.lexer.pos
	for {
		tokenPrec, offset, ok, binop := p.checkBinOpPrec()
		if !ok {
			tokenPrec = -1
		}
		if tokenPrec < expressionPrec {
			return lhs
		}

		for i := 0; i < offset; i++ {
			p.lexer.nextToken()
		}
		rhs := p.parseUnary()
		if rhs == nil {
			return nil
		}

		nextPrec, _, ok, _ := p.checkBinOpPrec()
		if !ok {
			nextPrec = -1
		}
		if tokenPrec < nextPrec {
			rhs = p.parseBinOpRHS(tokenPrec+1, rhs)
			if rhs == nil {
				return nil
			}
		}

		lhs = &BinaryAST{
			pos,
			astBinary,
			binop,
			lhs,
			rhs,
		}
	}
}

func (p *Parser) parsePrimary() AST {
	switch p.lexer.token {
	case TokIdentifier:
		return p.parseIdentifier()
	//case TokStr:
	//	return p.parseStr()
	case TokNumber:
		return p.parseNumber()
	case TokLParen:
		return p.parseParen()
	case TokTrue, TokFalse:
		return p.parseBool()
	default:
		what := tokens[p.lexer.token]
		if p.lexer.token == TokUnknown {
			what = string(p.lexer.unknownVal)
		}
		p.addError("Expression '"+what+"' has been used incorrectly.")
		return nil
	}
}

func (p *Parser) parseStmt() AST {
	switch p.lexer.token {
	case TokIf:
		return p.parseIfElse()
	case TokReturn:
		return p.parseReturn()
	//case TokForLoop:
	//	return p.parseLoop()
	default:
		what := tokens[p.lexer.token]
		if p.lexer.token == TokUnknown {
			what = string(p.lexer.unknownVal)
		}
		p.addError("Statement '"+what+"' has been used incorrectly.")
		return nil
	}
}

func (p *Parser) parseBool() AST {
	pos := p.lexer.pos
	val := 0

	if p.lexer.identifier == "true" {
		val = 1
	} else if p.lexer.identifier != "false" {
		p.addError("Error occurred while parsing boolean")
	}

	p.lexer.nextToken()
	return &BoolAST{pos, astBool, val}
}

func (p *Parser) parseParen() AST {
	p.lexer.nextToken()
	val := p.parseExpression()
	if val == nil {
		return nil
	}

	if p.lexer.token != TokRParen {
		p.addError("Parenthesis are not closed.")
	}

	p.lexer.nextToken()
	return val
}

func (p *Parser) parseIdentifier() AST {
	pos := p.lexer.pos
	name := p.lexer.identifier

	p.lexer.nextToken()

	if p.lexer.token != TokLParen {
		varType, ok := p.knownVars[name]
		if !ok {
			p.addError(fmt.Sprintf(`Variable "%s" does not exist!`, name))
		}

		return &VariableAST{
			Pos: 	  pos,
			kind:     astVariable,
			Name:     name,
			VarType:  varType,
		}
	}

	p.lexer.nextToken()

	var args []AST

	for p.lexer.token != TokRParen {
		if p.lexer.token == TokEOF {
			p.addError("Function call is not closed")
		}

		arg := p.parseExpression()
		if arg != nil {
			args = append(args, arg)
		}

		if p.lexer.token == TokRParen {
			break
		}

		if p.lexer.token != TokArgSep {
			p.addError("Expected ',' in '"+name+"' function call.")
		}

		p.lexer.nextToken()
	}

	p.lexer.nextToken()
	return &CallAST{pos, astCall, name, args}
}

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

func (p *Parser) parseIfElse() AST {
	pos := p.lexer.pos
	p.lexer.nextToken()

	cond := p.parseExpression()
	if cond == nil {
		p.addError("Syntax Error: No condition inside if")
	}
	scopePos := p.checkAndNext(TokLBrace)

	var trueBody []AST
	for p.lexer.token != TokRBrace {
		if p.lexer.token == TokEOF {
			p.addError("No closing brace in 'if' statement")
		}

		body := p.parseStmt()
		if body != nil {
			trueBody = append(trueBody, body)
		}
	}

	p.lexer.nextToken()

	var elseScope Pos
	var elseBody []AST
	var elseIfBody []ElseIfAST
	for {
		if p.lexer.token == TokEOF {
			p.addError("No closing brace in 'else if' statement")
		}

		if p.lexer.token != TokElse {
			break
		}

		p.lexer.nextToken()

		if p.lexer.token != TokIf {
			elseScope = p.checkAndNext(TokLBrace)
			for p.lexer.token != TokRBrace {
				if p.lexer.token == TokEOF {
					p.addError("No closing brace in 'else' statement")
				}

				body := p.parseStmt()
				if body != nil {
					elseBody = append(elseBody, body)
				}
			}

			p.lexer.nextToken()
			break
		}

		p.lexer.nextToken()

		elseIfCond := p.parseExpression()
		if elseIfCond == nil {
			p.addError("No condition inside 'else if'")
		}

		elseIfScope := p.checkAndNext(TokLBrace)

		var tempBody []AST
		for p.lexer.token != TokRBrace {
			if p.lexer.token == TokEOF {
				p.addError("No closing brace in 'else if' statement")
			}

			body := p.parseStmt()
			if body != nil {
				tempBody = append(tempBody, body)
			}
		}

		p.lexer.nextToken()
		elseIfBody = append(elseIfBody, ElseIfAST{
			elseIfScope,
			astIfElse,
			elseIfCond,
			BlockAST{
				elseIfScope,
				astBlock,
				tempBody,
			},
		})
	}

	return &IfElseAST{
		pos,
		astIfElse,
		cond,
		BlockAST{
			scopePos,
			astBlock,
			trueBody,
		},
		BlockAST{
			elseScope,
			astBlock,
			elseBody,
		},
		elseIfBody,
	}
}

func (p *Parser) parseReturn() AST {
	pos := p.lexer.pos
	p.lexer.ignoreNewLine = false
	p.lexer.nextToken()
	p.lexer.ignoreNewLine = true

	if p.lexer.token == TokUnknown && (p.lexer.unknownVal == 10 || p.lexer.unknownVal == 13) {
		return &ReturnAST{
			pos,
			astReturn,
			nil,
		}
	}

	value := p.parseExpression()
	return &ReturnAST{
		pos,
		astReturn,
		value,
	}
}

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

func (p *Parser) parseAssign(panicMessage string) interface{} {
	if p.lexer.token != TokAssign {
		panic(panicMessage)
	}

	p.lexer.nextToken()
	if p.lexer.token == TokIdentifier || p.lexer.token == TokAtom {
		return p.lexer.identifier
	}

	if p.lexer.token == TokNumber {
		return p.lexer.numVal
	}

	panic(panicMessage)
}

func (p *Parser) parsePrimitiveAttr() {
	p.lexer.nextToken()
	_ = p.checkAndNext(TokLParen)

	var prevToken Token
	for prevToken != TokRParen {
		if p.lexer.token == TokRParen {
			if prevToken == TokArgSep {
				p.addError("Wrong attribute definition. Excepted: ','")
			}

			break
		}

		if p.lexer.token == TokEOF {
			p.addError("Primitive attribute is not closed")
		}

		if p.lexer.token != TokIdentifier {
			p.addError("No identifier in the primitive attribute")
		}

		switch p.lexer.identifier {
		case "type":
			p.lexer.nextToken()
			typ := p.parseAssign("Invalid value assigning in the 'type' option of the primitive attribute")
			switch typ {
			case ":unary":
				p.isBinaryOp = false
			case ":binary":
				p.isBinaryOp = true
			default:
				p.addError("Type '" + typ.(string) + "' in the primitive attribute does not exist")
			}
		case "precedence":
			p.lexer.nextToken()
			precedence, ok := p.parseAssign("Invalid value assigning in the 'precedence' option of the primitive attribute").(float64)
			if !ok {
				p.addError("Could not assign value to precedence because value is not a number")
			}

			p.defaultPrecedence = int(precedence)
		default:
			p.addError("There is no '" + p.lexer.identifier + "' option in the primitive attribute")
		}

		p.lexer.nextToken()
		if p.lexer.token != TokArgSep && p.lexer.token != TokRParen {
			p.addError("Wrong attribute definition. Expected ',' or ')'.")
		}

		prevToken = p.lexer.token
		p.lexer.nextToken()
	}

	p.isOperator = true
}

func (p *Parser) parseAttribute() {
	p.lexer.nextToken()
	for p.lexer.unknownVal != ']' {
		if p.lexer.token == TokEOF {
			p.addError("Attribute is not closed")
		}

		if p.lexer.token != TokIdentifier {
			p.addError("Syntax Error: No identifier in the attribute")
		}

		switch p.lexer.identifier {
		case "primitive":
			p.parsePrimitiveAttr()
		default:
			p.addError("Attribute Error: '" + p.lexer.identifier + "' does not exist")
		}

		if p.lexer.token != TokArgSep && p.lexer.unknownVal != ']' {
			p.addError("Wrong attribute definition. Expected: ',' or ']'")
		}
	}

	p.lexer.nextToken()
}

func (p *Parser) parseUnary() AST {
	pos := p.lexer.pos
	if p.lexer.token != TokUnknown || p.lexer.unknownVal == ' ' {
		return p.parsePrimary()
	}

	unaryOp := p.lexer.unknownVal
	p.lexer.nextToken()
	if op := p.parseUnary(); op != nil {
		return &UnaryAST{
			Pos: 	  pos,
			kind:     astUnary,
			Operator: int(unaryOp),
			Operand:  op,
		}
	}

	return nil
}