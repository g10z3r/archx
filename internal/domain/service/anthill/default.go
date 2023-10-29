package anthill

import (
	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/common"
	"github.com/g10z3r/archx/internal/domain/service/anthill/config"
)

type CompassOption func(*Compass)

func DefaultCompass(options ...CompassOption) *Compass {
	compass := &Compass{
		config: &config.Config{
			RootDir:   ".",
			TargetDir: "",
			Analysis:  make(common.AnalyzerMap),
		},

		eventCh:       make(chan compassEvent, 1),
		unsubscribeCh: make(chan struct{}),
	}

	compass.RegisterAnalyzer(&analyzer.ImportAnalyzer{})
	compass.RegisterAnalyzer(&analyzer.StructAnalyzer{})
	compass.RegisterAnalyzer(&analyzer.FunctionAnalyzer{})

	for _, opt := range options {
		opt(compass)
	}

	return compass
}

func WithRootDir(dir string) CompassOption {
	return func(c *Compass) {
		c.config.RootDir = dir
	}
}

func WithTargetDir(dir string) CompassOption {
	return func(c *Compass) {
		c.config.TargetDir = dir
	}
}
