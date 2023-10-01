package scanner

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
	mongodbScanDAO "github.com/g10z3r/archx/internal/infrastructure/db/mongodb/scanner/dao"
)

type structRepository struct {
	documentID primitive.ObjectID
	collection *mongo.Collection
}

func newStructRepository(docID primitive.ObjectID, col *mongo.Collection) *structRepository {
	return &structRepository{
		documentID: docID,
		collection: col,
	}
}

func (r *structRepository) Append(ctx context.Context, structDTO *domainDTO.StructDTO, structIndex int, pkgPath string) error {
	filter := bson.D{
		{Key: "_id", Value: r.documentID},
		{Key: "packages.path", Value: pkgPath},
	}

	update := bson.D{
		{
			Key: "$push", Value: bson.D{
				{Key: "packages.$.structs", Value: mongodbScanDAO.MapStructDTO(structDTO)},
			},
		},
		{
			Key: "$set", Value: bson.D{
				{Key: fmt.Sprintf("packages.$.structsIndex.%s", *structDTO.Name), Value: structIndex},
			},
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
