package scaner

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path"

	"github.com/g10z3r/archx/internal/scaner/entity"
	"github.com/g10z3r/archx/pkg/bloom"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Scanner struct {
	cache      *scannerCache
	mongoDbCol *mongo.Collection
	documentID primitive.ObjectID
}

func NewScanner(mongoDbCol *mongo.Collection, documentID primitive.ObjectID) *Scanner {
	return &Scanner{
		cache:      newScannerCache(),
		mongoDbCol: mongoDbCol,
		documentID: documentID,
	}
}

func (s *Scanner) ScanPackage(client *mongo.Client, dirPath string, mod string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		newPackage := Package{
			Path:    dirPath,
			Name:    pkg.Name,
			Structs: []entity.Struct{},
			Imports: make([]string, 0),
		}

		filter := bson.D{
			{Key: "_id", Value: s.documentID},
		}

		update := bson.D{
			{Key: "$push", Value: bson.D{
				{Key: "packages", Value: newPackage},
			}},

			{Key: "$set", Value: bson.D{
				{Key: "packagesIndex." + newPackage.Name, Value: len(s.cache.packagesIndex)},
			}},
		}

		updateResult, err := s.mongoDbCol.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
		s.cache.AddPackage(newPackage.Path, len(s.cache.packagesIndex))

		pkgImports, total := processImports(pkg.Files)
		pkgCache := newPackageCache(mod, bloom.FilterConfig{
			ExpectedItemCount:        uint64(total),
			DesiredFalsePositiveRate: 0.01,
		})
		for _, _import := range pkgImports {
			contains, err := pkgCache.importsFilter.MightContain([]byte(_import.Path))
			if err != nil {
				log.Fatal(err)
			}

			if !contains {
				update := bson.D{
					{Key: "$addToSet", Value: bson.D{
						{Key: "packages.$[pkg].imports", Value: _import.Path},
					}},
				}

				// Определите опции массива
				arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
					Filters: []interface{}{bson.D{{Key: "pkg.path", Value: newPackage.Path}}},
				})

				updateResult, err := s.mongoDbCol.UpdateOne(context.Background(), filter, update, arrayFilters)
				if err != nil {
					log.Fatal(err)
				}

				pkgCache.AddImport(getAlias(_import), len(pkgCache.importsIndex))

				log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

			} else {
				if index := pkgCache.GetImportIndex(getAlias(_import)); index < 0 {
					pkgCache.AddImport(getAlias(_import), len(pkgCache.importsIndex))
				}
			}
		}

		log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}
}

func getAlias(imp *entity.Import) string {
	if imp.WithAlias {
		return imp.Alias
	}

	return path.Base(imp.Path)
}

func isSideEffectImport(imp *entity.Import) bool {
	return imp.WithAlias && imp.Alias == "_"
}

func processImports(files map[string]*ast.File) ([]*entity.Import, int) {
	var impTotal int
	var imports []*entity.Import

	for _, file := range files {
		impTotal = impTotal + len(file.Imports)

		for _, imp := range file.Imports {
			if imp.Path != nil && imp.Path.Value != "" {
				imports = append(imports, entity.NewImport(imp))
			}
		}
	}

	return imports, impTotal
}

// func (s *Scanner) ScanPackage2(client *mongo.Client, dirPath string, mod string) (*buffer.BufferEventBus, error) {
// 	var buf *buffer.BufferEventBus

// 	fset := token.NewFileSet()
// 	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, pkg := range pkgs {
// 		imports, total := processImports(pkg.Files)
// 		buf = buffer.NewBufferEventBus(mod, total, errChan)

// 		var wg sync.WaitGroup
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			buf.Open()
// 		}()

// 		for i := 0; i < len(imports); i++ {
// 			buf.SendEvent(&buffer.AddImportEvent{Import: imports[i]})
// 		}

// 		newPackage := Package{
// 			Path:    dirPath,
// 			Structs: []entity.Struct{},
// 		}

// 		filter := bson.D{
// 			{Key: "_id", Value: insertResult.InsertedID},
// 		}

// 		update := bson.D{
// 			{Key: "$push", Value: bson.D{
// 				{Key: "packages", Value: newPackage},
// 			}},
// 		}
// 		updateResult, err := collection.UpdateOne(context.Background(), filter, update)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

// 		for fileName, file := range pkg.Files {
// 			log.Printf("Processing file: %s", fileName)

// 			buf.WaitGroup.Add(1)
// 			go func(file *ast.File) {
// 				defer buf.WaitGroup.Done()

// 				for _, decl := range file.Decls {
// 					switch d := decl.(type) {
// 					case *ast.FuncDecl:
// 						processFuncDecl(buf, fset, d)
// 					case *ast.GenDecl:
// 						processGenDecl(collection, insertResult.InsertedID, buf, fset, d)
// 					}
// 				}
// 			}(file)
// 		}

// 		buf.WaitGroup.Wait()
// 		buf.Close()
// 		wg.Wait()
// 	}

// 	return buf, nil
// }
