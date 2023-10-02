package scanner

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/g10z3r/archx/internal/domain/repository"

	"github.com/g10z3r/archx/internal/domain/entity"
	"github.com/g10z3r/archx/internal/infrastructure/db/mongodb/scanner/model"
)

type packageRepository struct {
	documentID primitive.ObjectID
	collection *mongo.Collection

	importRepo repository.ImportRepository
	structRepo repository.StructRepository
}

func newPackageRepository(docID primitive.ObjectID, col *mongo.Collection) *packageRepository {
	return &packageRepository{
		documentID: docID,
		collection: col,

		importRepo: newImportRepository(docID, col),
		structRepo: newStructRepository(docID, col),
	}
}

func (r *packageRepository) ImportRepo() repository.ImportRepository {
	return r.importRepo
}

func (r *packageRepository) StructRepo() repository.StructRepository {
	return r.structRepo
}

func (r *packageRepository) Append(ctx context.Context, newPackage *entity.PackageEntity, packageIndex int) error {
	filter := bson.D{
		{Key: "_id", Value: r.documentID},
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "packages", Value: model.MapPackageEntity(newPackage)},
		}},

		{Key: "$set", Value: bson.D{
			{Key: "packagesIndex." + newPackage.Path, Value: packageIndex},
		}},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
