package config

import "github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"

type Config struct {
	Analysis   analyzer.AnalyzerMap
	ModuleName string
	RootDir    string
	TargetDir  string
}
