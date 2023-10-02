package main

import (
	_ "github.com/g10z3r/archx/example/internal/api"
	"github.com/g10z3r/archx/example/internal/data"
	apiWithAlias "github.com/g10z3r/archx/example/internal/data/api"
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

type Address struct {
	street, city string
	state        string
	zip          string
}

type Person struct {
	FirstName,
	LastName string
	Age     int
	test    apiWithAlias.Api2
	Address Address
	Info    data.PersonalInfo
	// Transport Transport
	Skill struct {
		Intelligence float32
		Info2        data.PersonalInfo
		Speed        float32
	}
}

// func (p *Person) ChangeFirstName(newFirstName string) *Person {
// 	pi := data.PersonalInfo{}
// 	pi.TestMethod()

// 	if newFirstName != p.FirstName {
// 		p.FirstName = newFirstName
// 	}

// 	return p
// }

// func (p *Person) ChangeLastName(newLastName string) *Person {
// 	os.Getenv("TEST")
// 	metadata.MetadataSome()
// 	if newLastName != p.LastName {
// 		p.LastName = newLastName
// 	}

// 	return p
// }
