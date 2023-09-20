package main

import (
	"github.com/g10z3r/archx/example/internal/api"
	mt "github.com/g10z3r/archx/example/internal/metadata"
)

type Test struct {
	Field1 string
	Field2 int
	Api    api.ApitestStruct
}

func (p *Person) AgeInc() *Person {
	mt.MetadataSome()
	p.Age = p.Age + 1

	return p
}
