package scanner

import (
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type scannerDb struct {
	mu sync.RWMutex

	collection *mongo.Collection
	documentID primitive.ObjectID
}

func newScannerDb(collection *mongo.Collection, documentID primitive.ObjectID) *scannerDb {
	return &scannerDb{
		mu: sync.RWMutex{},
	}
}
