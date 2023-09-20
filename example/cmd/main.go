package main

import (
	"os"

	"github.com/g10z3r/archx/example/internal/data"
)

type Go interface {
	Straight() error
	Left() error
	Right() error
	Back() error
}

type Transport interface {
	StartEngine() error
	StopEngine() error

	Go
}

type Address struct {
	street string
	city   string
	state  string
	zip    string
}

type Person struct {
	FirstName,
	LastName string
	Age       int
	Info      data.PersonalInfo
	Address   Address
	Transport Transport
	Skill     struct {
		Intelligence float32

		Speed float32
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

func (p *Person) ChangeLastName(newLastName string) *Person {
	os.Getenv("TEST")

	if newLastName != p.LastName {
		p.LastName = newLastName
	}

	return p
}
