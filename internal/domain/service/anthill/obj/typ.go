package obj

type AstTyp uint

const (
	ImportSpec AstTyp = iota + 1
	TypeScpec
	FuncDecl
	StructType
	FuncType
	InterfaceType
)
