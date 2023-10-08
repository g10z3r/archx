package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/g10z3r/archx/internal/domain/entity"
	"github.com/g10z3r/archx/internal/domain/service/anthill"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://ant:password@localhost:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	// db := client.Database("archant")
	// collection := db.Collection("someproject")

	// scanRepo := mongoScannerRepo.NewSnapshotRepository(collection)
	// scanService.Perform(ctx, "example/cmd", "github.com/g10z3r/archx")

	colony := anthill.SpawnColony(anthill.DefaultConfig(
		anthill.WithSelectedDir("example/cmd"),
	))

	if err := colony.Explore("."); err != nil {
		log.Fatal(err)
	}

	snapshot := entity.NewSnapshotEntity(colony.Metadata.ModName, len(colony.Packages))
	fmt.Println(colony.Packages)
	for _, pkg := range colony.Packages {
		ent, err := colony.Forage(pkg)
		if err != nil {
			log.Fatal(err)
		}

		snapshot.Packages = append(snapshot.Packages, ent)
	}

	jsonData, _ := json.Marshal(snapshot)
	fmt.Println(string(jsonData))
}
