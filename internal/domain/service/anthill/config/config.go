package config

import "github.com/g10z3r/archx/internal/domain/service/anthill/common"

type Config struct {
	Analysis   common.AnalyzerMap
	ModuleName string
	RootDir    string
	TargetDir  string
}
