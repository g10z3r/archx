package anthill

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/g10z3r/archx/internal/domain/entity"
	"golang.org/x/mod/modfile"
)

const goModFileName = "go.mod"
const goFileExt = ".go"

type Metadata struct {
	modName   string
	goVersion string
}

type Colony struct {
	mu sync.Mutex

	snapshot *entity.SnapshotEntity
	packages []string

	metadata *Metadata
	config   *Config
}

func SpawnColony(cfg *Config) *Colony {
	return &Colony{
		config:   cfg,
		metadata: &Metadata{},
	}
}

func (c *Colony) Explore(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	subdirs, goFilesExist, err := c.scanDirectory(entries, root)
	if err != nil {
		return err
	}

	if goFilesExist {
		if !strings.HasPrefix(root, c.config.selectedDir) {
			return nil
		}

		c.packages = append(c.packages, root)
		return nil
	}

	return c.exploreSubdirectories(subdirs)
}

func (c *Colony) scanDirectory(entries []os.DirEntry, root string) ([]string, bool, error) {
	var subdirs []string
	var goFilesExist bool

	for _, entry := range entries {
		if entry.IsDir() {
			subdir := filepath.Join(root, entry.Name())
			if _, exists := c.config.ignoredList[entry.Name()]; !exists {
				subdirs = append(subdirs, subdir)
			}
			continue
		}

		if entry.Name() == goModFileName && len(c.metadata.goVersion) < 1 {
			if err := c.processGoModFile(root); err != nil {
				return nil, false, err
			}
		}

		if filepath.Ext(entry.Name()) == goFileExt {
			goFilesExist = true
			break
		}
	}

	return subdirs, goFilesExist, nil
}

func (c *Colony) exploreSubdirectories(subdirs []string) error {
	for _, subdir := range subdirs {
		if err := c.Explore(subdir); err != nil {
			return err
		}
	}

	return nil
}

func (c *Colony) processGoModFile(root string) error {
	goModPath := filepath.Join(root, goModFileName)
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	modFileData, err := modfile.Parse(goModFileName, content, nil)
	if err != nil {
		return err
	}

	c.metadata.modName = modFileData.Module.Mod.Path
	c.metadata.goVersion = modFileData.Go.Version

	return nil
}

func (c *Colony) Forage() (*entity.SnapshotEntity, error) {
	if len(c.metadata.modName) < 1 || len(c.metadata.goVersion) < 1 {
		return nil, errors.New("couldn't find the go.mod file")
	}

	c.snapshot = entity.NewSnapshotEntity(c.metadata.modName, len(c.packages))

	for _, dir := range c.packages {
		if err := c.scanPackages(dir); err != nil {
			log.Fatal(err)
		}
	}

	return c.snapshot, nil
}

func (c *Colony) scanPackages(dirPath string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, pkg := range pkgs {
		wg.Add(1)
		go func(pkg *ast.Package) {
			forager := newForager(fset)
			pkgEntity := forager.process(pkg, dirPath, c.metadata.modName)

			c.mu.Lock()
			c.snapshot.Packages = append(c.snapshot.Packages, pkgEntity)
			c.mu.Unlock()

			wg.Done()
		}(pkg)
	}

	wg.Wait()

	return nil
}
