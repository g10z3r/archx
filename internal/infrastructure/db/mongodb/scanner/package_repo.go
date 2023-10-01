package scanner

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/g10z3r/archx/internal/domain/repository"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
	mongodbScanDAO "github.com/g10z3r/archx/internal/infrastructure/db/mongodb/scanner/dao"
)

type packageRepository struct {
	documentID primitive.ObjectID
	collection *mongo.Collection

	importRepo repository.ImportRepository
}

func newPackageRepository(docID primitive.ObjectID, col *mongo.Collection) *packageRepository {
	return &packageRepository{
		documentID: docID,
		collection: col,

		importRepo: newImportRepository(docID, col),
	}
}

func (r *packageRepository) ImportRepo() repository.ImportRepository {
	return r.importRepo
}

func (r *packageRepository) Append(ctx context.Context, newPackage *domainDTO.PackageDTO, packageIndex int) error {
	filter := bson.D{
		{Key: "_id", Value: r.documentID},
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "packages", Value: mongodbScanDAO.MapPackageDTO(newPackage)},
		}},

		{Key: "$set", Value: bson.D{
			{Key: "packagesIndex." + newPackage.Name, Value: packageIndex},
		}},
	}

	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
