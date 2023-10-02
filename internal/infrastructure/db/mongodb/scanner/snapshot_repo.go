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
	documentID  primitive.ObjectID
	collection  *mongo.Collection
	packageRepo repository.PackageRepository
}

func NewSnapshotRepository(col *mongo.Collection) repository.SnapshotRepository {
	return &snapshotRepository{collection: col}
}

func (r *snapshotRepository) PackageRepo() repository.PackageRepository {
	return r.packageRepo
}

func (r *snapshotRepository) Register(ctx context.Context, result *entity.SnapshotEntity) error {
	insertResult, err := r.collection.InsertOne(ctx, model.MapSnapshotEntity(result))
	if err != nil {
		return err
	}

	r.documentID = insertResult.InsertedID.(primitive.ObjectID)
	r.packageRepo = newPackageRepository(r.documentID, r.collection)

	return nil
}
