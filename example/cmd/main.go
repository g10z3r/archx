package main

import (
	"os"

	_ "github.com/g10z3r/archx/example/internal/api"
	"github.com/g10z3r/archx/example/internal/data"
	api2WithAnotherAlias "github.com/g10z3r/archx/example/internal/data/api"
	"github.com/g10z3r/archx/example/internal/metadata"
)

// type Go interface {
// 	Straight() error
// 	Left() error
// 	Right() error
// 	Back() error
// }

// type Transport interface {
// 	StartEngine() error
// 	StopEngine() error

// 	Go
// }

type SomeFunc[G comparable] func(a int, b string, g G) (uint, error)

type Address struct {
	street, city string
	state        string
	zip          string
}

type Person struct {
	FirstName,
	LastName string

	Age     int
	test    api2WithAnotherAlias.Api2
	Address Address
	Info    data.PersonalInfo
	Skill   struct {
		Intelligence float32
		Info2        data.PersonalInfo
		Speed        float32
	}
}

func (p *Person) ChangeFirstName(newFirstName string) *Person {
	pi := data.PersonalInfo{}
	pi.TestMethod()

	if newFirstName != p.FirstName {
		p.FirstName = newFirstName
	}

	return p
}

func (p *Person) ChangeLastName(newLastName string, d data.PersonalInfo) *Person {
	os.Getenv("")
	metadata.MetadataSome()
	if newLastName != p.LastName {
		p.LastName = newLastName
	}

	return p
}

func (p *Person) RecursiveMethod(data string) *Person {
	return p.RecursiveMethod(data)
}

func recursive(data string) *Person {
	return recursive(data)
}
