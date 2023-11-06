package main

import (
	"github.com/g10z3r/archx/example/internal/api"
	api2WithAnotherAlias "github.com/g10z3r/archx/example/internal/data"
	api2 "github.com/g10z3r/archx/example/internal/data/api"
)

type (
	AstructTest[A api2.GenrecicInterface, B any] struct {
		Field1 api2WithAnotherAlias.PersonalInfo
		Field2 api2.Api2
		Api    api.ApitestStruct
	}

	Bstruct struct {
		Filed3 string
		Field4 int
	}
)
