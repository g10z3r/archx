package scanner

import (
	"context"

	"github.com/g10z3r/archx/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
	mongodbScanDAO "github.com/g10z3r/archx/internal/infrastructure/db/mongodb/scanner/dao"
)

type scanResultRepository struct {
	documentID  primitive.ObjectID
	collection  *mongo.Collection
	packageRepo repository.PackageRepository
}

func NewScanResultRepository(col *mongo.Collection) repository.ScannerRepository {
	return &scanResultRepository{collection: col}
}

func (r *scanResultRepository) Create(ctx context.Context, result *domainDTO.ScanResultDTO) error {
	insertResult, err := r.collection.InsertOne(ctx, mongodbScanDAO.MapScanResultDTO(result))
	if err != nil {
		return err
	}

	r.documentID = insertResult.InsertedID.(primitive.ObjectID)
	r.packageRepo = newPackageRepository(r.documentID, r.collection)

	return nil
}

func (r *scanResultRepository) PackageRepo() repository.PackageRepository {
	return r.packageRepo
}
