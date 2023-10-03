package main

import (
	"github.com/g10z3r/archx/example/internal/api"
	api2WithAnotherAlias "github.com/g10z3r/archx/example/internal/data"
	api2 "github.com/g10z3r/archx/example/internal/data/api"
)

type Test struct {
	Field1 api2WithAnotherAlias.PersonalInfo
	Field2 api2.Api2
	Api    api.ApitestStruct
}

// func (p *Person) AgeInc() *Person {
// 	// _ = api.ApitestStruct{}
// 	p.Age = p.Age + 1

// 	return p
// }
