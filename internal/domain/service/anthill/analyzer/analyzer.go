package analyzer

import "go/ast"

type Object interface {
	Type() string
}

type Context interface{}

type Analyzer interface {
	Name() string
	Check(node ast.Node) bool
	Analyze(ctx *VisitorContext, spec ast.Node) Object
}
