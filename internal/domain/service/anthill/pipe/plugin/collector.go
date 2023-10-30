package plugin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

const goModFileName = "go.mod"
const goFileExt = ".go"

type CollectorPluginInput struct {
	IgnoredList map[string]struct{}
	PackageDirs []string
	RootDir     string
	TargetDir   string
}

type CollectorPluginOutput struct {
	PackageDirs []string
	ModuleName  string
	LangVersion string
}

type CollectorPlugin struct {
	next   Plugin
	name   string
	input  *CollectorPluginInput
	output *CollectorPluginOutput
}

func (p *CollectorPlugin) IsTerminal() bool {
	return false
}

func (p *CollectorPlugin) Next() Plugin {
	return p.next
}

func (p *CollectorPlugin) SetNext(plugin Plugin) {
	p.next = plugin
}

func (p *CollectorPlugin) Execute(ctx context.Context, input interface{}) interface{} {
	p.input = input.(*CollectorPluginInput)

	if err := p.explore(p.input.RootDir, true); err != nil {
		log.Fatal(err)
	}

	return p.input.PackageDirs
}

func (p *CollectorPlugin) explore(path string, isRoot bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	subdirs, goFilesExist, err := p.scanDir(entries, path)
	if err != nil {
		return err
	}

	if goFilesExist && strings.HasPrefix(path, p.input.TargetDir) {
		p.input.PackageDirs = append(p.input.PackageDirs, path)
	}

	if err := p.exploreSubDir(subdirs); err != nil {
		return err
	}

	if !isRoot {
		return nil
	}

	if len(p.output.ModuleName) < 1 || len(p.output.LangVersion) < 1 {
		return errors.New("couldn't find the go.mod file")
	}

	return nil
}

func (c *CollectorPlugin) scanDir(entries []os.DirEntry, root string) ([]string, bool, error) {
	var subdirs []string
	goFilesExist := false

	for _, entry := range entries {
		entryName := entry.Name()

		if entry.IsDir() {
			if _, exists := c.input.IgnoredList[entryName]; exists {
				continue
			}
			subdirs = append(subdirs, filepath.Join(root, entryName))
			continue
		}

		if entryName == goModFileName {
			if err := c.processGoMod(root); err != nil {
				return nil, false, fmt.Errorf("failed to process go.mod: %w", err)
			}
			continue
		}

		if filepath.Ext(entryName) == goFileExt {
			goFilesExist = true
		}
	}

	return subdirs, goFilesExist, nil
}

func (c *CollectorPlugin) exploreSubDir(subdirs []string) error {
	for _, subdir := range subdirs {
		if err := c.explore(subdir, false); err != nil {
			return err
		}
	}

	return nil
}

func (c *CollectorPlugin) processGoMod(root string) error {
	goModPath := filepath.Join(root, goModFileName)
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	modFileData, err := modfile.Parse(goModFileName, content, nil)
	if err != nil {
		return err
	}

	c.output.ModuleName = modFileData.Module.Mod.Path
	c.output.LangVersion = modFileData.Go.Version

	return nil
}
