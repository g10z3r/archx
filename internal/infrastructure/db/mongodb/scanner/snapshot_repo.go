package scanner

import (
	"context"

	"github.com/g10z3r/archx/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/g10z3r/archx/internal/domain/entity"
	"github.com/g10z3r/archx/internal/infrastructure/db/mongodb/scanner/model"
)

type snapshotRepository struct {
	documentID primitive.ObjectID
	collection *mongo.Collection

	packageAcc repository.PackageAccessor
}

func NewSnapshotRepository(col *mongo.Collection) repository.SnapshotRepository {
	return &snapshotRepository{collection: col}
}

func (r *snapshotRepository) PackageAcc() repository.PackageAccessor {
	return r.packageAcc
}

func (r *snapshotRepository) Register(ctx context.Context, result *entity.SnapshotEntity) error {
	insertResult, err := r.collection.InsertOne(ctx, model.MapSnapshotEntity(result))
	if err != nil {
		return err
	}

	r.documentID = insertResult.InsertedID.(primitive.ObjectID)
	r.packageAcc = newPackageAccessor(r.documentID, r.collection)

	return nil
}
