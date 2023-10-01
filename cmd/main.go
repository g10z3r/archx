package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/g10z3r/archx/internal/scaner"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// total := 100
	// var ts string
	// var counter int
	// for i := 0; i < total; i++ {
	// 	buf, err := scaner.ScanPackage("./example/cmd", "github.com/g10z3r/archx")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	jsonData, err := json.Marshal(buf)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	fmt.Println(string(jsonData))

	// 	if ts != "" && len(ts) != len(string(jsonData)) {
	// 		fmt.Println(i)
	// 		break
	// 	}

	// 	ts = string(jsonData)
	// 	counter++

	// }

	// if counter == total {
	// 	fmt.Println("\n\nDone!")
	// }

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

	db := client.Database("archant")
	collection := db.Collection("someproject")
	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"packages.path": 1},
		Options: options.Index().SetUnique(true),
	}
	indexName, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created index: %s\n", indexName)

	insertResult, err := collection.InsertOne(context.Background(), scaner.ScanResult{
		Timestamp:     time.Now().Unix(),
		Packages:      []scaner.Package{},
		PackagesIndex: make(map[string]int),
	})
	if err != nil {
		log.Fatal(err)
	}

	// s := scaner.Scanner{
	// 	MongoDbCol: collection,
	// 	DocumentID: insertResult.InsertedID.(primitive.ObjectID),
	// }

	s := scaner.NewScanner(
		collection,
		insertResult.InsertedID.(primitive.ObjectID),
	)

	s.ScanPackage(client, "./example/cmd", "github.com/g10z3r/archx")

}
