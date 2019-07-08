package main

import (
	"ExpectoLang/llvm/bindings/go/llvm"
	"fmt"
	"os"
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
	fcPassManager.AddLoopUnswitchPass()
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
	case "+":
		return builder.CreateFAdd(l, r, "addtmp")
	case "-":
		return builder.CreateFSub(l, r, "subtmp")
	case "*":
		return builder.CreateFMul(l, r, "multmo")
	case "/":
		// TODO check if r is not 0
		return builder.CreateFDiv(l, r, "divtmp")
	case "<":
		return builder.CreateFCmp(llvm.FloatOLT, l, r, "cmptmp")
	case "==":
		return builder.CreateFCmp(llvm.FloatOEQ, l, r, "cmptmp")
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
	case "+":
		// Todo figure out string concat
		panic("not implemented: String Concat")
	case "==":
		l = builder.CreatePointerCast(l, llvm.Int8Type(), "pointcast")
		r = builder.CreatePointerCast(r, llvm.Int8Type(), "pointcast")
		return builder.CreateICmp(llvm.IntEQ, l, r, "cmptmp")
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
		panic(fmt.Sprintf(`Function "%s" could not be referenced`, c.Callee))
	}

	if callee.ParamsCount() != len(c.args) {
		panic(fmt.Sprintf(`Incorrect arguments passed in function "%s"`, c.Callee))
	}

	argsValues := []llvm.Value{}

	for _, arg := range c.args {
		argVal := arg.codegen()
		if argVal.IsNil() {
			panic(fmt.Sprintf(`One of the arguments in function "%s" was null`, c.Callee))
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
		// TODO THIS IS MY DEBUG
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

	if os.Getenv("DEBUG") != "true" {
		fcPassManager.RunFunc(proc)
	}

	return proc
}

func (i *IfElseAST) codegen() llvm.Value {
	cond := i.Condition.codegen()
	if cond.IsNil() {
		panic("No condition")
	}

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

func (r *ReturnAST) codegen() llvm.Value {
	if r.Body == nil {
		return builder.CreateRetVoid()
	}
	return builder.CreateRet(r.Body.codegen())
}

func (l *LoopAST) codegen() llvm.Value {
	cond := l.Condition.codegen()
	if cond.IsNil() {
		panic("No condition in the loop")
	}
	zeroInd := llvm.ConstInt(llvm.Int32Type(), 0, false)

	elemAlloca := builder.CreateArrayAlloca(cond.Operand(0).Operand(0).Type().ElementType(), llvm.ConstInt(llvm.Int32Type(), 1, false), "")
	gep := llvm.ConstGEP(cond, []llvm.Value{zeroInd})
	condAlloca := builder.CreateAlloca(gep.Type(), "")

	builder.CreateStore(gep, condAlloca)
	load := builder.CreateLoad(condAlloca, "load")
	gep2 := builder.CreateInBoundsGEP(load, []llvm.Value{zeroInd}, "")
	load2 := builder.CreateLoad(gep2, "load")
	builder.CreateStore(load2, elemAlloca)

	fc := builder.GetInsertBlock().Parent()
	headerBlock := builder.GetInsertBlock()
	loopBlock := llvm.AddBasicBlock(fc, "loop")
	builder.CreateBr(loopBlock)

	builder.SetInsertPointAtEnd(loopBlock)
	valInd := builder.CreatePHI(llvm.Int32Type(), "ind")
	valInd.AddIncoming([]llvm.Value{llvm.ConstInt(llvm.Int32Type(), 0, false)}, []llvm.BasicBlock{headerBlock})

	oldValInd, okInd := namedValues[l.IndexVar]
	namedValues[l.IndexVar] = valInd

	oldValElem, okElem := namedValues[l.ElementVar]
	namedValues[l.ElementVar] = elemAlloca

	// TODO Check if loop's body does not have return inside
	l.Body.codegen()

	// Get next element from array and get next index
	nextInd := builder.CreateAdd(valInd, llvm.ConstInt(llvm.Int32Type(), 1, false), "nextind")
	gep2 = builder.CreateInBoundsGEP(load, []llvm.Value{nextInd}, "")
	load2 = builder.CreateLoad(gep2, "load")
	builder.CreateStore(load2, elemAlloca)

	breakCond := builder.CreateICmp(llvm.IntNE, llvm.ConstInt(llvm.Int32Type(), uint64(cond.Operand(0).Operand(0).Type().ArrayLength()), false), nextInd, "loopcond")

	loopExitBlock := builder.GetInsertBlock()
	exitBlock := llvm.AddBasicBlock(fc, "exitloop")
	builder.CreateCondBr(breakCond, loopBlock, exitBlock)
	builder.SetInsertPointAtEnd(exitBlock)
	valInd.AddIncoming([]llvm.Value{nextInd}, []llvm.BasicBlock{loopExitBlock})

	if okInd {
		namedValues[l.IndexVar] = oldValInd
	} else {
		delete(namedValues, l.IndexVar)
	}

	if okElem {
		namedValues[l.ElementVar] = oldValElem
	} else {
		delete(namedValues, l.ElementVar)
	}

	return llvm.ConstNull(llvm.Int1Type())
}