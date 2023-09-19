package entity

import (
	"go/ast"
	"go/token"
)

type InterfaceType struct {
	_   [0]int
	pos token.Pos
	end token.Pos
}

func NewInterfaceType(res *ast.InterfaceType) *InterfaceType {
	return &InterfaceType{
		pos: res.Pos(),
		end: res.End(),
	}
}
