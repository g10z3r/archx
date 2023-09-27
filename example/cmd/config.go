package main

import (
	"github.com/g10z3r/archx/example/internal/api"
	api2 "github.com/g10z3r/archx/example/internal/data/api"
)

type Test struct {
	Field1 string
	Field2 api2.Api2
	Api    api.ApitestStruct
}
