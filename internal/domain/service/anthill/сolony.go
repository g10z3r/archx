package anthill

import (
	"context"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

const goModFileName = "go.mod"
const goFileExt = ".go"

type Colony struct {
	zone      []string
	mod       string
	goVersion string
	config    Config
}

func NewColony(cfg Config) *Colony {
	return &Colony{
		config: cfg,
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
		c.zone = append(c.zone, root)
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

		if entry.Name() == goModFileName && len(c.goVersion) < 1 {
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

func (s *Colony) Perform(ctx context.Context, dirPath string, basePath string) {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {

		thread := newForager(fset)
		thread.process(pkg, basePath)

	}
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

	c.mod = modFileData.Module.Mod.Path
	c.goVersion = modFileData.Go.Version

	return nil
}

func (c *Colony) processEntries() {

}
