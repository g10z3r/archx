package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/g10z3r/archx/internal/domain/service/anthill"
	"github.com/g10z3r/archx/internal/domain/service/anthill/collector"
	"github.com/g10z3r/archx/internal/domain/service/anthill/event"
)

var ignoredMap = map[string]struct{}{
	".git":    {},
	".docker": {},

	".vscode":  {},
	".idea":    {},
	".eclipse": {},

	"dist":    {},
	"docker":  {},
	"assets":  {},
	"vendor":  {},
	"build":   {},
	"scripts": {},
	"ci":      {},
	"log":     {},
	"logs":    {},
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

	compass := anthill.NewCompass()
	clct := collector.DefaultCollector(
		collector.WithTargetDir("example/cmd"),
	)
	if err := clct.Explore(); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	eventCh, unsubscribeCh := compass.Subscribe()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case e := <-eventCh:
				switch ev := e.(type) {
				case *event.PackageFormedEvent:
					jsonData, err := json.Marshal(ev.Package)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(string(jsonData))

				default:
					fmt.Printf("Unknown event type: %s\n", e.Name())
				}
			case <-unsubscribeCh:
				return
			}
		}
	}()

	for _, p := range clct.GetAllPackageDirs() {
		compass.Parse(clct.GetInfo(), p)
		time.Sleep(time.Second)
		close(unsubscribeCh)
	}
}
