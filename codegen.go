package main

import (
	"github.com/llir/ll/ast"
	"github.com/llir/llvm/ir"
)

var (
	module = ir.NewModule()
	namedValues = map[string]ast.Value{}
)

func (n *NumberLiteralAST) codegen() ast.Value {
	return ast.FloatConst{

	}
}