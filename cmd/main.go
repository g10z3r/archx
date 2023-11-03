package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"log"
	"log/slog"

	"github.com/g10z3r/archx/internal/domain/service/anthill"
	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/collector"
)

type MyWriter struct {
	data []byte
}

// Write реализует метод интерфейса io.Writer
func (w *MyWriter) Write(p []byte) (int, error) {
	// Здесь можно определить, как обрабатывать переданные данные
	w.data = append(w.data, p...)
	return len(p), nil
}

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// clientOptions := options.Client().ApplyURI("mongodb://ant:password@localhost:27017")

	// client, err := mongo.Connect(ctx, clientOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = client.Ping(ctx, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Successfully connected to MongoDB!")

	// db := client.Database("archant")
	// collection := db.Collection("someproject")

	// scanRepo := mongoScannerRepo.NewSnapshotRepository(collection)
	// scanService.Perform(ctx, "example/cmd", "github.com/g10z3r/archx")

	myWriter := &MyWriter{}

	logger := slog.New(slog.NewTextHandler(myWriter, nil))
	fmt.Println(len(myWriter.data))
	logger.Info("Test")
	logger.Info("Test1")

	fmt.Printf("Записанные логи:\n%s\n", myWriter.data)

	// compass.Run(context.Background())
	clct := collector.DefaultCollector(
		collector.WithTargetDir("example"),
	)
	if err := clct.Explore(); err != nil {
		log.Fatal(err)
	}

	compass := anthill.NewEngine(&anthill.EngineConfig{
		AnalyzerFactoryMap: getAnalyzers(),
		Determinator:       baseNodeDeterminator,
		ModuleName:         clct.GetInfo().ModuleName,
	})

	// var wg sync.WaitGroup

	// eventCh := make(chan event.Event)
	// unsubscribeCh := compass.Subscribe(eventCh)

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for {
	// 		select {
	// 		case e := <-eventCh:
	// 			switch ev := e.(type) {
	// 			case *event.PackageFormedEvent:
	// 				jsonData, err := json.Marshal(ev.Package)
	// 				if err != nil {
	// 					log.Fatal(err)
	// 				}

	// 				fmt.Println(string(jsonData))

	// 			default:
	// 				fmt.Printf("Unknown event type: %s\n", e.Name())
	// 			}
	// 		case <-unsubscribeCh:
	// 			return
	// 		}
	// 	}
	// }()

	for _, p := range clct.GetAllPackageDirs() {
		data, err := compass.ParseDir(p)
		if err != nil {

		}

		for _, pkg := range data {
			jsonData, err := json.Marshal(pkg)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("\n", string(jsonData))
		}

	}
}

// TODO: tmp func
func getAnalyzers() anthill.EngineAFMap {
	return anthill.EngineAFMap{
		anthill.ImportNodeType: analyzer.NewImportAnalyzer,
		anthill.FuncNodeType:   analyzer.NewFuncAnalyzer,
		anthill.StructNodeType: analyzer.NewStructAnalyzer,
	}
}

func baseNodeDeterminator(node ast.Node) uint {
	switch n := node.(type) {
	case *ast.ImportSpec:
		return anthill.ImportNodeType
	case *ast.FuncDecl:
		return anthill.FuncNodeType
	// case *ast.FuncType:
	// 	return anthill.FuncNodeType
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.StructType:
			return anthill.StructNodeType
		}
	}

	return 0
}

// func getAnalyzers() anthill.EngineAFMap {
// 	return anthill.EngineAFMap{
// 		reflect.TypeOf(new(ast.ImportSpec)): analyzer.NewImportAnalyzer,
// 		reflect.TypeOf(new(ast.FuncDecl)):   analyzer.NewFuncAnalyzer,
// 		reflect.TypeOf(new(ast.TypeSpec)):   analyzer.NewStructAnalyzer,
// 	}
// }
