package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/g10z3r/archx/internal/domain/service/anthill"
	"github.com/g10z3r/archx/internal/domain/service/anthill/event"
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type SnapshotEntity struct {
	Timestamp int64
	BasePath  string
	Packages  []*obj.PackageObj
}

func NewSnapshotEntity(mod string, pkgCount int) *SnapshotEntity {
	return &SnapshotEntity{
		Timestamp: time.Now().Unix(),
		BasePath:  mod,
	}
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

	// colony := anthill.SpawnColony(anthill.DefaultConfig(
	// 	anthill.WithSelectedDir("example/cmd"),
	// ))

	compass := anthill.NewCompass()

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

	compass.Parse()
	time.Sleep(time.Second)
	close(unsubscribeCh)
	wg.Wait()
	// if err := colony.Explore("."); err != nil {
	// 	log.Fatal(err)
	// }

	// snapshot := NewSnapshotEntity(colony.Metadata.ModName, len(colony.Packages))
	// fmt.Println(colony.Packages)
	// for _, pkg := range colony.Packages {
	// 	ent, err := colony.Forage(pkg)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	snapshot.Packages = append(snapshot.Packages, ent)
	// }
}
