package anthill

import (
	"go/ast"
	"go/token"

	"github.com/g10z3r/archx/internal/domain/entity"
)

func (f *forager) processGenDecl(fset *token.FileSet, genDecl *ast.GenDecl, impMeta map[string]int, fileName string) error {
	if genDecl.Tok != token.TYPE {
		return nil
	}

	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		switch t := typeSpec.Type.(type) {
		case *ast.StructType:
			structEntity, err := f.processStructType(fset, newStructProcDTO(typeSpec, t), impMeta)
			if err != nil {
				return err
			}

			bucket, _ := f.storage.structBucket.Load(fileName)
			f.storage.structBucket.Store(fileName, append(bucket, structEntity))
		}
	}

	return nil
}

type structProcDTO struct {
	typeSpec   *ast.TypeSpec
	structType *ast.StructType
	structName string
}

func newStructProcDTO(typeSpec *ast.TypeSpec, structType *ast.StructType) *structProcDTO {
	return &structProcDTO{
		typeSpec:   typeSpec,
		structType: structType,
		structName: typeSpec.Name.Name,
	}
}

func (f *forager) processStructType(fset *token.FileSet, dto *structProcDTO, impMeta map[string]int) (*entity.StructEntity, error) {
	structEntity, usedPackages, err := entity.NewStructEntity(fset, dto.structType, entity.NotEmbedded, &dto.structName)
	if err != nil {
		return nil, err
	}

	for _, pkg := range usedPackages {
		if index, exists := impMeta[pkg.Alias]; exists {
			structEntity.AddDependency(index, pkg.Element)
		}
	}

	// Synchronizing methods from the buffer that have already been found and are awaiting their structure
	// for _, method := range th.buf.GetAndClearMethods(params.structName) {
	// 	structEntity.AddMethod(method, method.Name)

	// 	depsLen := len(structEntity.DependenciesIndex)
	// 	for dep, i := range method.DependenciesIndex {
	// 		// When a dependency already exists in structure dependencies,
	// 		// just increase the usage counter for this dependency
	// 		if j, exists := structEntity.DependenciesIndex[dep]; exists {
	// 			structEntity.Dependencies[j].Usage++
	// 			continue
	// 		}

	// 		// This is a new dependency that has never been seen before
	// 		structEntity.Dependencies = append(structEntity.Dependencies, method.Dependencies[i])
	// 		structEntity.DependenciesIndex[dep] = depsLen + i
	// 	}
	// }

	return structEntity, nil
}
