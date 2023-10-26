package collector

import "github.com/g10z3r/archx/internal/domain/service/anthill/config"

type CollectorOption func(*Collector)

func WithIgnoredList(list ...string) CollectorOption {
	return func(c *Collector) {
		for _, item := range list {
			c.ignoredList[item] = struct{}{}
		}
	}
}

func WithTargetDir(dir string) CollectorOption {
	return func(c *Collector) {
		c.targetDir = dir
	}
}

func WithRootDir(dir string) CollectorOption {
	return func(c *Collector) {
		c.rootDir = dir
	}
}

func DefaultCollector(options ...CollectorOption) *Collector {
	collector := &Collector{
		ignoredList: config.DefaultIgnoredMap,
		rootDir:     ".",
		targetDir:   "",
		projInfo:    &ProjectInfo{},
		Packages:    make([]string, 0),
	}

	for _, opt := range options {
		opt(collector)
	}

	return collector
}
