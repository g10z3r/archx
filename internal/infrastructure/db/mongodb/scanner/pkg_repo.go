package scanner

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
	mongoDBScanDAO "github.com/g10z3r/archx/internal/infrastructure/db/mongodb/scanner/dao"
)

type packageRepository struct {
	documentID primitive.ObjectID
	collection *mongo.Collection
}

func newPackageRepository(docID primitive.ObjectID, collection *mongo.Collection) *packageRepository {
	return &packageRepository{
		documentID: docID,
		collection: collection,
	}
}

func (r *packageRepository) Create(ctx context.Context, newPackage *domainDTO.PackageDTO, packageIndex int) error {
	filter := bson.D{
		{Key: "_id", Value: r.documentID},
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "packages", Value: mongoDBScanDAO.MapPackageDTO(newPackage)},
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
