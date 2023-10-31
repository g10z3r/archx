package scanner

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type importAccessor struct {
	documentID primitive.ObjectID
	collection *mongo.Collection
}

func newImportAccessor(docID primitive.ObjectID, col *mongo.Collection) *importAccessor {
	return &importAccessor{
		documentID: docID,
		collection: col,
	}
}

func (r *importAccessor) Append(ctx context.Context, _import *obj.ImportObj, packagePath string) error {
	filter := bson.D{
		{Key: "_id", Value: r.documentID},
	}

	update := bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "packages.$[pkg].imports", Value: _import.Path},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.D{{Key: "pkg.path", Value: packagePath}}},
	})

	_, err := r.collection.UpdateOne(ctx, filter, update, arrayFilters)
	return err
}

func (r *importAccessor) AppendSideEffectImport(ctx context.Context, _import *obj.ImportObj, packagePath string) error {
	filter := bson.D{
		{Key: "_id", Value: r.documentID},
	}

	update := bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "packages.$[pkg].sideEffectImports", Value: _import.Path},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.D{{Key: "pkg.path", Value: packagePath}}},
	})

	_, err := r.collection.UpdateOne(ctx, filter, update, arrayFilters)
	return err
}
