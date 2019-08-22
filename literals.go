package main

import "strconv"

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


