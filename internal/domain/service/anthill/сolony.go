package anthill

import (
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/g10z3r/archx/internal/domain/obj"
	"golang.org/x/mod/modfile"
)

const goModFileName = "go.mod"
const goFileExt = ".go"

type Metadata struct {
	ModName   string
	GoVersion string
}

type Colony struct {
	mu       sync.Mutex
	Packages []string
	Metadata *Metadata
	config   *Config
}

func SpawnColony(cfg *Config) *Colony {
	return &Colony{
		config:   cfg,
		Metadata: &Metadata{},
	}
}

func (c *Colony) Explore(root string) error {
	return c.recursiveExplore(root, true)
}

func (c *Colony) recursiveExplore(root string, isRootCall bool) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	subdirs, goFilesExist, err := c.scanDirectory(entries, root)
	if err != nil {
		return err
	}

	if goFilesExist && strings.HasPrefix(root, c.config.selectedDir) {
		c.Packages = append(c.Packages, root)
	}

	if err := c.exploreSubdirectories(subdirs); err != nil {
		return err
	}

	if !isRootCall {
		return nil
	}

	if len(c.Metadata.ModName) < 1 || len(c.Metadata.GoVersion) < 1 {
		return errors.New("couldn't find the go.mod file")
	}

	return nil
}

func (c *Colony) scanDirectory(entries []os.DirEntry, root string) ([]string, bool, error) {
	var subdirs []string
	goFilesExist := false

	for _, entry := range entries {
		entryName := entry.Name()

		if entry.IsDir() {
			if _, exists := c.config.ignoredList[entryName]; exists {
				continue
			}
			subdirs = append(subdirs, filepath.Join(root, entryName))
			continue
		}

		if entryName == goModFileName && len(c.Metadata.GoVersion) < 1 {
			if err := c.processGoModFile(root); err != nil {
				return nil, false, err
			}
			continue
		}

		if filepath.Ext(entryName) == goFileExt {
			goFilesExist = true
		}
	}

	return subdirs, goFilesExist, nil
}

func (c *Colony) exploreSubdirectories(subdirs []string) error {
	for _, subdir := range subdirs {
		if err := c.recursiveExplore(subdir, false); err != nil {
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

	c.Metadata.ModName = modFileData.Module.Mod.Path
	c.Metadata.GoVersion = modFileData.Go.Version

	return nil
}

func (c *Colony) Forage(dirPath string) (*obj.PackageObj, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	var pkgEntity *obj.PackageObj
	forager := newForager(fset)

	for _, pkg := range pkgs {
		pkgEntity = forager.process(pkg, dirPath, c.Metadata.ModName)
		break
	}

	return pkgEntity, nil
}
