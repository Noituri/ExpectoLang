package main

import (
	"ExpectoLang/llvm/bindings/go/llvm"
	"fmt"
	"strings"
)

var (
	module        llvm.Module
	builder       = llvm.NewBuilder()
	namedValues   = map[string]llvm.Value{}
	fcPassManager llvm.PassManager
)

func InitModuleAndPassManager() {
	module = llvm.NewModule("expectoroot")
	fcPassManager = llvm.NewFunctionPassManagerForModule(module)
	fcPassManager.AddInstructionCombiningPass()
	fcPassManager.AddGVNPass()
	fcPassManager.AddCFGSimplificationPass()
	fcPassManager.InitializeFunc()
}

func (s *StringAST) codegen() llvm.Value {
	return builder.CreateGlobalStringPtr(strings.ReplaceAll(s.Value, `\n`, "\n"), "strtmp")
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

func (b *BinaryAST) numberCodegen() llvm.Value {
	l := b.Lhs.codegen()
	r := b.Rhs.codegen()

	if l.IsNil() || r.IsNil() {
		panic("null operands")
	}

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

func (b *BinaryAST) strCodegen() llvm.Value {
	l := b.Lhs.codegen()
	r := b.Rhs.codegen()

	if l.IsNil() || r.IsNil() {
		panic("null operands")
	}

	switch b.Op {
	case '+':
		// Todo figure out string concat
		panic("not implemented: String Concat")
	default:
		panic(fmt.Sprintf(`Operator "%c" is invalid`, b.Op))
	}
}

func (b *BinaryAST) codegen() llvm.Value {
	if b.Lhs.Kind() == astNumber && b.Rhs.Kind() == astNumber {
		return b.numberCodegen()
	} else if b.Lhs.Kind() == astString && b.Rhs.Kind() == astString {
		return b.strCodegen()
	}

	return b.numberCodegen()

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

		argsValues = append(argsValues, argVal)
	}

	return builder.CreateCall(callee, argsValues, "")
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
	var fcType llvm.Type

	switch p.ReturnType {
	case LitFloat:
		fcType = llvm.FunctionType(llvm.FloatType(), args, false)
	case LitString:
		fcType = llvm.FunctionType(llvm.PointerType(llvm.Int8Type(), 0), args, false)
	case LitVoid:
		fcType = llvm.FunctionType(llvm.VoidType(), args, false)
	default:
		// TODO DEBUG
		//fcType = llvm.FunctionType(llvm.FloatType(), args, false)
		//fcType = llvm.FunctionType(llvm.Int32Type(), args, false)
		panic(fmt.Sprintf("type-%s-does-no-exit", p.ReturnType))
	}

	fc := llvm.AddFunction(module, p.Name, fcType)

	for i, param := range fc.Params() {
		param.SetName(p.Args[i].Name)
	}

	return fc
}

func (b *BlockAST) codegen() ([]llvm.Value, bool) {
	elements := []llvm.Value{}
	isReturn := false
	for _, stmt := range b.Elements {
		elements = append(elements, stmt.codegen())
		if stmt.Kind() == astReturn {
			isReturn = true
			break
		}
	}
	return elements, isReturn
}

// TODO Check for redefinition
// TODO Check if function has return
func (p *FunctionAST) codegen() llvm.Value {
	proc := module.NamedFunction(p.Proto.Name)

	if proc.IsNil() {
		proc = p.Proto.codegen()
	}

	if proc.IsNil() {
		panic(fmt.Sprintf(`Could not create function "%s"`, p.Proto.Name))
	}
	block := llvm.AddBasicBlock(proc, "entry")
	builder.SetInsertPointAtEnd(block)

	namedValues = map[string]llvm.Value{}

	for _, param := range proc.Params() {
		namedValues[param.Name()] = param
	}

	p.Body.codegen()

	if llvm.VerifyFunction(proc, llvm.PrintMessageAction) != nil {
		proc.EraseFromParentAsFunction()
		panic(fmt.Sprintf(`Error occurred while verifing function "%s"`, p.Proto.Name))
	}

	fcPassManager.RunFunc(proc)

	return proc
}

func (i *IfElseAST) codegen() llvm.Value {
	cond := i.Condition.codegen()
	if cond.IsNil() {
		panic("No condition")
	}

	cond = builder.CreateFCmp(llvm.FloatONE, cond, llvm.ConstFloat(llvm.FloatType(), 0), "cond")

	fc := builder.GetInsertBlock().Parent()
	thenBlock := llvm.AddBasicBlock(fc, "then")
	elseBlock := llvm.AddBasicBlock(fc, "else")
	exitBlock := llvm.AddBasicBlock(fc, "exit")

	builder.CreateCondBr(cond, thenBlock, elseBlock)

	// build then body
	builder.SetInsertPointAtEnd(thenBlock)
	_, isRet := i.TrueBody.codegen()

	if !isRet {
		builder.CreateBr(exitBlock)
	}

	// build elifs body
	for ind, el := range i.ElifBody {
		if ind == 0 {
			builder.SetInsertPointAtEnd(elseBlock)
			elseBlock = llvm.AddBasicBlock(fc, "else")
		}

		elifThenBlock := llvm.AddBasicBlock(fc, "then")
		elifElseBlock := llvm.AddBasicBlock(fc, "else")

		elifCond := el.Condition.codegen()
		if elifCond.IsNil() {
			panic("No condition in elif")
		}

		elifCond = builder.CreateFCmp(llvm.FloatONE, elifCond, llvm.ConstFloat(llvm.FloatType(), 0), "cond")
		if ind == len(i.ElifBody) - 1 {
			builder.CreateCondBr(elifCond, elifThenBlock, elseBlock)
		} else {
			builder.CreateCondBr(elifCond, elifThenBlock, elifElseBlock)
		}

		builder.SetInsertPointAtEnd(elifThenBlock)
		_, isRet := el.Body.codegen()

		if !isRet {
			builder.CreateBr(exitBlock)
		}

		builder.SetInsertPointAtEnd(elifElseBlock)
		if ind == len(i.ElifBody) - 1 {
			builder.CreateBr(exitBlock)
		}
	}

	// build else body
	builder.SetInsertPointAtEnd(elseBlock)
	_, isRet = i.FalseBody.codegen()

	if !isRet {
		builder.CreateBr(exitBlock)
	}

	builder.SetInsertPointAtEnd(exitBlock)

	return cond
}

// TODO return might return void
func (r *ReturnAST) codegen() llvm.Value {
	if r.Body == nil {
		return builder.CreateRetVoid()
	}
	return builder.CreateRet(r.Body.codegen())
}