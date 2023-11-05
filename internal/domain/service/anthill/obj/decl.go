package obj

import (
	"go/ast"
	"go/token"
)

// Obtained after traversing the array of `*ast.Field`
type extractedFieldsData struct {
	usedPackages []UsedPackage
	fieldsSet    []*StructFieldObj
}

type FieldObj struct {
	Name string
	Type string
}

type DeclObj struct {
	Start token.Pos
	End   token.Pos
	Obj   Object
	// TODO: Introduce the use of internal object markings.
	// Stop using reflection in groups of anatomyizers and link everything to internal types of AST
	Typ AstTyp
}

func (o *DeclObj) Type() string {
	return "decl"
}

func NewDeclObj(node ast.Node, obj Object) *DeclObj {
	return &DeclObj{
		Start: node.Pos(),
		End:   node.End(),
		Typ:   0,
		Obj:   obj,
	}
}
