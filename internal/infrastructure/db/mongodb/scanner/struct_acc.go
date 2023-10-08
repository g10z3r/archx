package scanner

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/g10z3r/archx/internal/domain/obj"
	"github.com/g10z3r/archx/internal/infrastructure/db/mongodb/scanner/model"
)

type structAccessor struct {
	mu sync.Mutex

	documentID primitive.ObjectID
	collection *mongo.Collection
}

func newStructAccessor(docID primitive.ObjectID, col *mongo.Collection) *structAccessor {
	return &structAccessor{
		mu:         sync.Mutex{},
		documentID: docID,
		collection: col,
	}
}

func (r *structAccessor) Append(ctx context.Context, structEntity *obj.StructObj, structIndex int, pkgPath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filter := bson.D{
		{Key: "_id", Value: r.documentID},
		{Key: "packages.path", Value: pkgPath},
	}

	update := bson.D{
		{
			Key: "$push", Value: bson.D{
				{Key: "packages.$.structs", Value: model.MapStructEntity(structEntity)},
			},
		},
		{
			Key: "$set", Value: bson.D{
				{Key: fmt.Sprintf("packages.$.structsIndex.%s", *structEntity.Name), Value: structIndex},
			},
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	log.Printf("Saved %s to db", *structEntity.Name)
	return err
}

// func (r *structAccessor) AddMethod(ctx context.Context, methodEntity *entity.MethodEntity, structIndex int, pkgPath string) error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	// filter := bson.D{
// 	// 	{Key: "_id", Value: r.documentID},
// 	// 	{Key: "packages.path", Value: pkgPath},
// 	// }

// 	// if structIndex == -1 {
// 	// 	newStruct := model.StructDAO{
// 	// 		Fields:            make([]*model.FieldDAO, 0),
// 	// 		FieldsIndex:       make(map[string]int),
// 	// 		Methods:           []*model.MethodDAO{model.MapMethodEntity(methodEntity)},
// 	// 		Dependencies:      make([]*model.DependencyDAO, 0, len(methodEntity.Dependencies)),
// 	// 		DependenciesIndex: make(map[string]int, len(methodEntity.Dependencies)),
// 	// 	}

// 	// 	// r.Append(ctx, newStruct, )
// 	// }

// 	// update := bson.D{
// 	// 	{
// 	// 		Key: "$push", Value: bson.D{
// 	// 			{
// 	// 				Key:   fmt.Sprintf("packages.$.structs.%d.methods", structIndex),
// 	// 				Value: model.MapMethodEntity(methodEntity),
// 	// 			},
// 	// 		},
// 	// 	},
// 	// }

// 	// _, err := r.collection.UpdateOne(ctx, filter, update)
// 	return nil
// }
