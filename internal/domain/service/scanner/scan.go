package scanner

import (
	"context"
	"go/parser"
	"go/token"
	"log"
	"time"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
	"github.com/g10z3r/archx/internal/domain/repository"
)

type ScanService struct {
	cache *scannerCache
	db    repository.ScannerRepository
}

func NewScanService(scanRepo repository.ScannerRepository) *ScanService {
	return &ScanService{
		cache: newScannerCache(),
		db:    scanRepo,
	}
}

func (s *ScanService) Perform(ctx context.Context, dirPath string, mod string) {
	scanResult := domainDTO.ScannerResultDTO{
		Timestamp: time.Now().Unix(),
	}
	if err := s.db.Create(ctx, &scanResult); err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		newPkg := domainDTO.PackageDTO{
			Path: dirPath,
			Name: pkg.Name,
		}

		err := s.db.PackageRepo().Create(ctx, &newPkg, len(s.cache.packagesIndex))
		if err != nil {
			log.Fatal(err)
		}

	}
}
