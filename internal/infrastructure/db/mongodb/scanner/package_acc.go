package scanner

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/g10z3r/archx/internal/domain/repository"
)

type packageAccessor struct {
	documentID primitive.ObjectID
	collection *mongo.Collection

	importAcc repository.ImportAccessor
	structAcc repository.StructAccessor
}

func newPackageAccessor(docID primitive.ObjectID, col *mongo.Collection) *packageAccessor {
	return &packageAccessor{
		documentID: docID,
		collection: col,

		// importAcc: newImportAccessor(docID, col),
		// structAcc: newStructAccessor(docID, col),
	}
}

func (r *packageAccessor) ImportAcc() repository.ImportAccessor {
	return r.importAcc
}

func (r *packageAccessor) StructAcc() repository.StructAccessor {
	return r.structAcc
}

// func (r *packageAccessor) Append(ctx context.Context, newPackage *obj.PackageObj, packageIndex int) error {
// 	filter := bson.D{
// 		{Key: "_id", Value: r.documentID},
// 	}

// 	update := bson.D{
// 		{Key: "$push", Value: bson.D{
// 			{Key: "packages", Value: model.MapPackageEntity(newPackage)},
// 		}},

// 		{Key: "$set", Value: bson.D{
// 			{Key: "packagesIndex." + newPackage.Path, Value: packageIndex},
// 		}},
// 	}

// 	_, err := r.collection.UpdateOne(ctx, filter, update)
// 	return err
// }
