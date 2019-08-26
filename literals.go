package main

import (
	"ExpectoLang/llvm/bindings/go/llvm"
	"strconv"
)

const (
	LitVoid   = "void"
	LitString = "str"
	LitFloat  = "float"
	LitBool	  = "bool"
	LitInt	  = "int"
)

func ASTTypeToLit(astType kind) string {
	switch astType {
	case astString:
		return LitString
	case astNumberFloat:
		return LitFloat
	case astBool:
		return LitBool
	case astNumberInt:
		return LitInt
	}

	panic("Type of index '" + strconv.Itoa(int(astType)) + "' does not exist as literal")
}

func LLVMTypeToLit(llvmType llvm.Type) string {
	switch llvmType {
	case llvm.PointerType(llvm.Int8Type(), 0):
		return LitString
	case llvm.DoubleType():
		return LitFloat
	case llvm.Int1Type():
		return LitBool
	case llvm.Int32Type():
		return LitInt
	}

	panic("Type of index '" + llvmType.String() + "' does not exist as literal")
}


