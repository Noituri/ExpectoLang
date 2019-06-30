package main

import (
	"ExpectoLang/llvm/bindings/go/llvm"
	"fmt"
)

var (
	module      = llvm.NewModule("expectoroot")
	builder     = llvm.NewBuilder()
	namedValues = map[string]llvm.Value{}
)

func (s *StringAST) codegen() llvm.Value {
	return llvm.ConstArray(llvm.Int8Type(), []llvm.Value{llvm.ConstInt(llvm.Int8Type(), 44, false)})
}

func (n *NumberLiteralAST) codegen() llvm.Value {
	return llvm.ConstFloat(llvm.FloatType(), n.Value)
}

func (v *VariableAST) codegen() llvm.Value {
	val, ok := namedValues[v.Name]

	if !ok {
		panic(fmt.Sprintf(`Variable "%s" does not exist!`, v.Name))
	}

	return val
}

func (b *BinaryAST) codegen() llvm.Value {
	l := b.Lhs.codegen()
	r := b.Rhs.codegen()

	if l.IsNil() || r.IsNil() {
		panic("null operands")
	}

	// TODO check types
	switch b.Op {
	case '+':
		return builder.CreateFAdd(l, r, "addtmp")
	case '-':
		return builder.CreateFSub(l, r, "subtmp")
	case '*':
		return builder.CreateFMul(l, r, "multmo")
	case '/':
		// TODO check if r is not 0
		return builder.CreateFDiv(l, r, "divtmp")
	case '<':
		builder.CreateFCmp(llvm.FloatOLT, l, r, "cmptmp")
		return builder.CreateUIToFP(l, llvm.FloatType(), "booltmp")
	default:
		panic(fmt.Sprintf(`Operator "%c" is invalid`, b.Op))
	}
}

func (c *CallAST) codegen() llvm.Value {
	callee := module.NamedFunction(c.Callee)

	if callee.IsNil() {
		panic(fmt.Sprintf(`Procedure "%s" could not be referenced`, c.Callee))
	}

	if callee.ParamsCount() != len(c.args) {
		panic(fmt.Sprintf(`Incorrect arguments passed in procedure "%s"`, c.Callee))
	}

	argsValues := []llvm.Value{}

	for _, arg := range c.args {
		argVal := arg.codegen()

		if argVal.IsNil() {
			panic(fmt.Sprintf(`One of the arguments in procedure "%s" was null`, c.Callee))
		}
	}

	return builder.CreateCall(callee, argsValues, "calltmp")
}

func (p *PrototypeAST) codegen() llvm.Value {
	args := make([]llvm.Type, 0, len(p.Args))

	for _, a := range p.Args {
		// TODO use args type
		switch a.ArgType {
		case LitFloat:
			args = append(args, llvm.FloatType())
		case LitString:
			args = append(args, llvm.PointerType(llvm.Int8Type(), 0))
		default:
			panic(fmt.Sprintf("type-%s-does-no-exit", a.ArgType))
		}
	}

	// TODO use p.returnType
	procType := llvm.FunctionType(llvm.FloatType(), args, false)
	proc := llvm.AddFunction(module, p.Name, procType)

	for i, param := range proc.Params() {
		param.SetName(p.Args[i].Name)
	}

	return proc
}

// TODO use this to generate body
func (b *BlockAST) codegen() llvm.Value {
	return llvm.Value{}
}

// TODO Check for redefinition
func (p *FunctionAST) codegen() llvm.Value {
	proc := module.NamedFunction(p.Proto.Name)

	if proc.IsNil() {
		proc = p.Proto.codegen()
	}

	if proc.IsNil() {
		panic(fmt.Sprintf(`Could not create procedure "%s"`, p.Proto.Name))
	}
	block := llvm.AddBasicBlock(proc, "entry")
	builder.SetInsertPointAtEnd(block)

	namedValues = map[string]llvm.Value{}

	for _, param := range proc.Params() {
		namedValues[param.Name()] = param
	}

	for _, stmt := range p.Body.Body[:len(p.Body.Body)-1] {
		stmt.codegen()
	}

	retVal := p.Body.Body[len(p.Body.Body)-1].codegen()

	if retVal.IsNil() {
		panic(fmt.Sprintf(`No return in procedure "%s"`, p.Proto.Name))
	}

	builder.CreateRet(retVal)

	if llvm.VerifyFunction(proc, llvm.PrintMessageAction) != nil {
		proc.EraseFromParentAsFunction()
		panic(fmt.Sprintf(`Error occurred while verifing procedure "%s"`, p.Proto.Name))
	}

	return proc
}
