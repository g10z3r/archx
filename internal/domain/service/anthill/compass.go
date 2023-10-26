package anthill

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

var alz = map[string]analyzer.Analyzer{}

type Manager struct {
	analyzers map[string]analyzer.Analyzer
}

func (m *Manager) Register(a analyzer.Analyzer) {
	alz[a.Name()] = a
}

type Compass struct {
	fset    *token.FileSet
	manager *Manager
}

func NewCompass() *Compass {
	return &Compass{
		fset: token.NewFileSet(),
		manager: &Manager{
			analyzers: map[string]analyzer.Analyzer{},
		},
	}
}

func (c *Compass) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ImportSpec:

	case *ast.TypeSpec:
		if structType, ok := n.Type.(*ast.StructType); ok {
			obj.NewStructObj(c.fset, structType, obj.NotEmbedded, &n.Name.Name)
		}

		fmt.Println(n.Name.Name)
		// Handle ident nodes

	case *ast.FuncDecl:
		fmt.Println(n.Name.Name)
		// Handle func declaration nodes

	// Add cases for other node types...

	default:
		return c
	}

	return c
}

func (c *Compass) Parse() *obj.PackageObj {
	importAlz := &analyzer.ImportAnalyzer{}
	c.manager.Register(importAlz)

	structAlz := &analyzer.StructAnalyzer{}
	c.manager.Register(structAlz)

	// funcAlz := &analyzer.FunctionAnalyzer{}
	// c.manager.Register(funcAlz)

	fset := token.NewFileSet()
	pkg, err := parser.ParseDir(fset, "./example/cmd", nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	packa := &obj.PackageObj{}
	for _, p := range pkg {
		for _, f := range p.Files {
			vis := analyzer.NewVisitor(fset, alz, f.Name.Name)
			ast.Walk(vis, f)

			for analyzerName := range alz {
				switch analyzerName {
				case "import":
					data := toPkgImports(vis.Unload(analyzerName))
					packa.Imports = append(packa.Imports, data...)
				case "struct":
					data := toPkgStructs(vis.Unload(analyzerName))
					packa.Structs = append(packa.Structs, data...)
				}
			}
		}
	}

	return packa
}

func toPkgStructs(data []analyzer.Object) []*obj.StructObj {
	dataOutput := make([]*obj.StructObj, 0, len(data))

	for _, structObj := range data {
		dataOutput = append(dataOutput, structObj.(*obj.StructObj))
	}

	return dataOutput
}

func toPkgImports(data []analyzer.Object) []string {
	dataOutput := make([]string, 0, len(data))
	for _, importObj := range data {
		dataOutput = append(dataOutput, importObj.(*obj.ImportObj).Path)
	}

	return dataOutput
}
